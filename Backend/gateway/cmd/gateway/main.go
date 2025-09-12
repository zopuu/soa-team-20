package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	tb "github.com/didip/tollbooth/v7"
	tbchi "github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zopuu/soa-team-20/Backend/gateway/internal/config"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/mw"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/proxy"
	followerspb "github.com/zopuu/soa-team-20/Backend/services/followers_service/proto/followerspb"
)

func main() {
	cfg := config.New()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(40 * time.Second))
	r.Use(mw.CORS(cfg.CorsOrigins).Handler)

	// Rate limit (50 req/sec/IP)
	limiter := tb.NewLimiter(50, nil)
	limiter.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})
	r.Use(tbchi.LimitHandler(limiter))

	jwtCfg := mw.JWTConfig{
		Secret:   []byte(cfg.JwtSecret),
		Issuer:   cfg.JwtIssuer,
		Audience: cfg.JwtAudience,
	}

	// ---- Proxies ----
	authProxy, _ := proxy.NewHTTPReverseProxy(proxy.Options{
		Target: cfg.AuthBase, StripPrefix: "", DialTimeout: cfg.DialTimeout, ProxyTimeout: cfg.ProxyTimeout,
	})
	stakeProxy, _ := proxy.NewHTTPReverseProxy(proxy.Options{
		Target: cfg.StakeBase, StripPrefix: "", DialTimeout: cfg.DialTimeout, ProxyTimeout: cfg.ProxyTimeout,
	})
	blogProxy, _ := proxy.NewHTTPReverseProxy(proxy.Options{
		Target: cfg.BlogBase, StripPrefix: "", DialTimeout: cfg.DialTimeout, ProxyTimeout: cfg.ProxyTimeout,
	})
	tourProxy, _ := proxy.NewHTTPReverseProxy(proxy.Options{
		Target: cfg.TourBase, StripPrefix: "", DialTimeout: cfg.DialTimeout, ProxyTimeout: cfg.ProxyTimeout,
	})

	logProxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Gateway forwarding request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		log.Printf("Forwarding to %s%s", cfg.StakeBase, r.URL.Path)
		stakeProxy.ServeHTTP(w, r) // forward to users service
	})

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public auth endpoints (allow anonymous): e.g., /api/auth/login, /api/auth/register, /api/auth/refresh
	r.Route("/api/auth", func(rr chi.Router) {
		rr.Use(mw.AuthOptional(jwtCfg)) // parse token if present, but don't require
		rr.Handle("/*", authProxy)
	})

	// Everything else requires JWT
	secure := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request secured with JWT middleware: %s", r.Header.Get("Authorization"))
			mw.AuthRequired(jwtCfg)(h).ServeHTTP(w, r)
		})
	}

	r.Group(func(pr chi.Router) {
		pr.Route("/api/users", func(rr chi.Router) { rr.Handle("/*", secure(logProxy)) })
		pr.Route("/blogs", func(rr chi.Router) { rr.Handle("/*", secure(blogProxy)) })
		pr.Route("/tours", func(rr chi.Router) { rr.Handle("/*", secure(tourProxy)) })
		pr.Route("/keyPoints", func(rr chi.Router) { rr.Handle("/*", secure(tourProxy)) })
		pr.Route("/simulator", func(rr chi.Router) { rr.Handle("/*", secure(tourProxy)) })
	})

	// ---------- gRPC placeholder (next sprint) ----------
	// Here you'd init gRPC clients and expose REST handlers that call them:
	// grpcConn, _ := grpc.Dial(cfg.StakeGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// stakeClient := pb.NewStakeholdersClient(grpcConn)
	// r.Get("/api/bff/profile", func(w http.ResponseWriter, r *http.Request) { ... stakeClient.GetProfile(ctx, req) ... })
	// ----------------------------------------------------

	// Connect to Followers gRPC service
	grpcConn, err := grpc.Dial(cfg.FollowersGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to followers service: %v", err)
	}
	defer grpcConn.Close()

	followersClient := followerspb.NewFollowersServiceClient(grpcConn)

	// REST -> gRPC mappings
	r.Route("/api/followers", func(rr chi.Router) {
		rr.Use(secure) // require JWT auth

		rr.Post("/follow", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				FollowerId string `json:"user_id"`
				FolloweeId string `json:"target_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp, err := followersClient.Follow(r.Context(), &followerspb.FollowRequest{
				FollowerId: req.FollowerId,
				FolloweeId: req.FolloweeId,
			})
			log.Printf("Connected: %v", resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(resp)
		})

		rr.Post("/unfollow", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				UserID   string `json:"user_id"`
				TargetID string `json:"target_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp, err := followersClient.Unfollow(r.Context(), &followerspb.FollowRequest{
				FollowerId: req.UserID,
				FolloweeId: req.TargetID,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(resp)
		})

		rr.Get("/{id}/following", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			resp, err := followersClient.GetFollowing(r.Context(), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})

		rr.Get("/{id}/followers", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			resp, err := followersClient.GetFollowers(r.Context(), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})

		rr.Get("/recommendations/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			resp, err := followersClient.GetRecommendations(r.Context(), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})
	})

	addr := ":" + cfg.Port
	log.Printf("API Gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

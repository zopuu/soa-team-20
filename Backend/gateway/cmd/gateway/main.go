package main

import (
	"log"
	"net/http"
	"time"

	tb "github.com/didip/tollbooth/v7"
	tbchi "github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/zopuu/soa-team-20/Backend/gateway/internal/config"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/mw"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/proxy"
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
		return mw.AuthRequired(jwtCfg)(h)
	}

	r.Group(func(pr chi.Router) {
		pr.Route("/api/users", func(rr chi.Router) { rr.Handle("/*", secure(stakeProxy)) })
		pr.Route("/blogs",         func(rr chi.Router) { rr.Handle("/*", secure(blogProxy)) })
		pr.Route("/tours",         func(rr chi.Router) { rr.Handle("/*", secure(tourProxy)) })
		pr.Route("/keyPoints",         func(rr chi.Router) { rr.Handle("/*", secure(tourProxy)) })
	})

	// ---------- gRPC placeholder (next sprint) ----------
	// Here you'd init gRPC clients and expose REST handlers that call them:
	// grpcConn, _ := grpc.Dial(cfg.StakeGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// stakeClient := pb.NewStakeholdersClient(grpcConn)
	// r.Get("/api/bff/profile", func(w http.ResponseWriter, r *http.Request) { ... stakeClient.GetProfile(ctx, req) ... })
	// ----------------------------------------------------

	addr := ":" + cfg.Port
	log.Printf("API Gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

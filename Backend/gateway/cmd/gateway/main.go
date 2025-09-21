package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"

	"go.uber.org/zap"
    "go.uber.org/zap/zapcore"

	tb "github.com/didip/tollbooth/v7"
	tbchi "github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zopuu/soa-team-20/Backend/gateway/internal/config"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/mw"
	"github.com/zopuu/soa-team-20/Backend/gateway/internal/proxy"
	followerspb "github.com/zopuu/soa-team-20/Backend/services/followers_service/proto/followerspb"
)

type ctxKey string
const traceKey ctxKey = "trace_id"

func getTraceID(r *http.Request) string {
	if v, ok := r.Context().Value(traceKey).(string); ok && v != "" { return v }
	return ""
}
// korelacione headere dodajemo na outgoing HTTP zahteve
func withCorrHeaders(h http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.Header.Set("X-Request-ID", middleware.GetReqID(r.Context()))
        r.Header.Set("X-Trace-Id", getTraceID(r))
        h.ServeHTTP(w, r)
    }
}
// gRPC kontekst sa korelacionim ID-jevima
func grpcCtxWithTrace(r *http.Request) context.Context {
    md := metadata.Pairs(
        "x-request-id", middleware.GetReqID(r.Context()),
        "x-trace-id",   getTraceID(r),
    )
    return metadata.NewOutgoingContext(r.Context(), md)
}

func main() {
	cfg := config.New()

	// ---- logger ----
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.LevelKey = "level"
	encCfg.MessageKey = "msg"

	lvl := zapcore.InfoLevel
	if v := os.Getenv("LOG_LEVEL"); v != "" { _ = lvl.UnmarshalText([]byte(v)) }

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(os.Stdout), lvl)
	logger := zap.New(core).With(zap.String("service", "gateway"))
	defer logger.Sync()


	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(40 * time.Second))
	r.Use(mw.CORS(cfg.CorsOrigins).Handler)

	

	// 4a) TRACE middleware (obezbeđuje X-Trace-Id i stavlja u context + response header)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			traceID := r.Header.Get("X-Trace-Id")
			if traceID == "" { traceID = reqID } // jednostavno pravilo
			w.Header().Set("X-Trace-Id", traceID)
			ctx := context.WithValue(r.Context(), traceKey, traceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// 4b) JSON access-log (zamenjuje middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.String("client_ip", r.RemoteAddr),
				zap.Int64("latency_ms", time.Since(start).Milliseconds()),
				zap.String("request_id", middleware.GetReqID(r.Context())),
				zap.String("trace_id", getTraceID(r)),
			)
		})
	})


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

	/*
	logProxy := withCorrHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// informativni log pre forward-a
		logger.Info("forward_http",
			zap.String("to", cfg.StakeBase),
			zap.String("path", r.URL.Path),
			zap.String("trace_id", getTraceID(r)),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)
		stakeProxy.ServeHTTP(w, r)
	}))
		*/


	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Public auth endpoints (allow anonymous): e.g., /api/auth/login, /api/auth/register, /api/auth/refresh
	r.Route("/api/auth", func(rr chi.Router) {
		rr.Use(mw.AuthOptional(jwtCfg)) // parse token if present, but don't require
		rr.Handle("/*", withCorrHeaders(authProxy))
	})

	// Everything else requires JWT
	secure := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("secure_route",
				zap.Bool("authz_present", r.Header.Get("Authorization") != ""),
				zap.String("trace_id", getTraceID(r)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			)
			mw.AuthRequired(jwtCfg)(h).ServeHTTP(w, r)
		})
	}


	r.Group(func(pr chi.Router) {
		pr.Route("/blogs",     func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(blogProxy))) })
		pr.Route("/tours",     func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(tourProxy))) })
		pr.Route("/keyPoints", func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(tourProxy))) })
		pr.Route("/simulator", func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(tourProxy))) })
		pr.Route("/api/admin", func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(authProxy))) })
		pr.Route("/api/users", func(rr chi.Router){ rr.Handle("/*", secure(withCorrHeaders(stakeProxy))) })
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

			resp, err := followersClient.Follow(grpcCtxWithTrace(r), &followerspb.FollowRequest{
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

			resp, err := followersClient.Unfollow(grpcCtxWithTrace(r), &followerspb.FollowRequest{
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
			resp, err := followersClient.GetFollowing(grpcCtxWithTrace(r), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})

		rr.Get("/{id}/followers", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			resp, err := followersClient.GetFollowers(grpcCtxWithTrace(r), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})

		rr.Get("/recommendations/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			resp, err := followersClient.GetRecommendations(grpcCtxWithTrace(r), &followerspb.UserRequest{UserId: id})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})
	})

	addr := ":" + cfg.Port
	logger.Info("gateway_listening", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("server_exit", zap.Error(err))
	}

}

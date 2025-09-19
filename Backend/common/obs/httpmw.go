package obs

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKey string

const (
	CtxReqID   ctxKey = "request_id"
	CtxTraceID ctxKey = "trace_id"
)

func ReqIDFrom(r *http.Request) string {
	if v, ok := r.Context().Value(CtxReqID).(string); ok {
		return v
	}
	return ""
}
func TraceIDFrom(r *http.Request) string {
	if v, ok := r.Context().Value(CtxTraceID).(string); ok {
		return v
	}
	return ""
}

// Dodaje/uzima X-Request-ID i X-Trace-Id; stavlja u ctx i response header.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = reqID
		}

		w.Header().Set("X-Request-ID", reqID)
		w.Header().Set("X-Trace-Id", traceID)

		ctx := context.WithValue(r.Context(), CtxReqID, reqID)
		ctx = context.WithValue(ctx, CtxTraceID, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type lrw struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *lrw) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *lrw) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return n, err
}

// JSON access-log (metod, path, status, bytes, ip, trajanje, req/trace id)
func AccessLogMiddleware(l *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &lrw{ResponseWriter: w}
			next.ServeHTTP(ww, r)

			l.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Int("bytes", ww.bytes),
				zap.String("client_ip", clientIP(r)),
				zap.Int64("latency_ms", time.Since(start).Milliseconds()),
				zap.String("request_id", ReqIDFrom(r)),
				zap.String("trace_id", TraceIDFrom(r)),
			)
		})
	}
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		return xf
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

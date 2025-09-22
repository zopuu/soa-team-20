package obs

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Metrics struct {
	svc   string
	reg   *prometheus.Registry
	httpC *prometheus.CounterVec
	httpH *prometheus.HistogramVec
	grpcC *prometheus.CounterVec
	grpcH *prometheus.HistogramVec
}

func NewMetrics(service string) *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	m := &Metrics{
		svc: service,
		reg: reg,
		httpC: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "HTTP request count",
			},
			[]string{"service", "method", "route", "code"},
		),
		httpH: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method", "route"},
		),
		grpcC: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "gRPC request count",
			},
			[]string{"service", "method", "code"},
		),
		grpcH: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "gRPC request duration",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
	}
	reg.MustRegister(m.httpC, m.httpH, m.grpcC, m.grpcH)
	return m
}

// ---------- HTTP (chi/mux/vanilla) ----------
type rwCounter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *rwCounter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }
func (w *rwCounter) Write(b []byte) (int, error) {
	if w.status == 0 { w.status = http.StatusOK }
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (m *Metrics) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &rwCounter{ResponseWriter: w}
		next.ServeHTTP(ww, r)

		route := r.URL.Path
		// If using gorilla/mux, prefer the templated route
		if cr := mux.CurrentRoute(r); cr != nil {
			if t, err := cr.GetPathTemplate(); err == nil && t != "" {
				route = t
			}
		}

		code := ww.status
		if code == 0 { code = http.StatusOK }

		m.httpC.WithLabelValues(m.svc, r.Method, route, strconv.Itoa(code)).Inc()
		m.httpH.WithLabelValues(m.svc, r.Method, route).Observe(time.Since(start).Seconds())
	})
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
}

// ---------- gRPC ----------
func (m *Metrics) GRPCUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := h(ctx, req)
		st, _ := status.FromError(err)
		code := "OK"
		if st != nil { code = st.Code().String() }

		m.grpcC.WithLabelValues(m.svc, info.FullMethod, code).Inc()
		m.grpcH.WithLabelValues(m.svc, info.FullMethod).Observe(time.Since(start).Seconds())
		return resp, err
	}
}

// Tiny HTTP server for metrics (useful for pure gRPC services).
func (m *Metrics) ServeMetrics(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.Handler())
	return http.ListenAndServe(addr, mux)
}

package obs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Helpers za čitanje iz context-a (korisno i van ovog paketa)
func ReqIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxReqID).(string); ok && v != "" {
		return v
	}
	return ""
}
func TraceIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(CtxTraceID).(string); ok && v != "" {
		return v
	}
	return ""
}

// Unary interceptor koji obezbeđuje X-Request-ID i X-Trace-Id u ctx + vraća ih u response headeru.
func GRPCTraceUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		var reqID, traceID string
		if md != nil {
			if v := md.Get("x-request-id"); len(v) > 0 {
				reqID = v[0]
			}
			if v := md.Get("x-trace-id"); len(v) > 0 {
				traceID = v[0]
			}
		}
		if reqID == "" {
			reqID = uuid.NewString()
		}
		if traceID == "" {
			traceID = reqID
		}

		// u ctx
		ctx = context.WithValue(ctx, CtxReqID, reqID)
		ctx = context.WithValue(ctx, CtxTraceID, traceID)

		// vrati klijentu kroz header
		_ = grpc.SetHeader(ctx, metadata.Pairs(
			"x-request-id", reqID,
			"x-trace-id", traceID,
		))

		return handler(ctx, req)
	}
}

// Access-log za unary gRPC pozive.
func GRPCAccessLogUnary(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)

		l.Info("grpc_request",
			zap.String("method", info.FullMethod),
			zap.String("grpc_code", st.Code().String()),
			zap.Int64("latency_ms", time.Since(start).Milliseconds()),
			zap.String("request_id", ReqIDFromCtx(ctx)),
			zap.String("trace_id", TraceIDFromCtx(ctx)),
		)
		return resp, err
	}
}

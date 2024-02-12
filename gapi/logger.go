package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	protocolKey   = "protocol"
	methodKey     = "method"
	pathKey       = "path"
	durationKey   = "duration"
	statusCodeKey = "statusCode"
	statusTextKey = "statusText"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str(protocolKey, "gRPC").
		Str(methodKey, info.FullMethod).
		Int(statusCodeKey, int(statusCode)).
		Str(statusTextKey, statusCode.String()).
		Dur(durationKey, duration).
		Msg("Received a gRPC request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{ResponseWriter: res, StatusCode: http.StatusOK}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)
		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str(protocolKey, "HTTP").
			Str(methodKey, req.Method).
			Str(pathKey, req.RequestURI).
			Int(statusCodeKey, rec.StatusCode).
			Str(statusTextKey, http.StatusText(rec.StatusCode)).
			Dur(durationKey, duration).
			Msg("Received a HTTP request")
	})
}

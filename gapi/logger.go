package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	protocolKey   = "protocol"
	methodKey     = "method"
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

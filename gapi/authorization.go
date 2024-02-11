package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicodanke/bankTutorial/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser (ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("Missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("Invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("Unsupported authorization type %s", authType)
	}
	// The JWT token is the second field in the authorization header
	tokenString := fields[1]

	payload, err := server.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("Inavalid access token: %s", err)
	}
	
	return payload, nil
}
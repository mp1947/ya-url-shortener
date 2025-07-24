// Package interceptor provides gRPC server interceptors for handling cross-cutting concerns
// such as authentication and authorization. These interceptors can be used to validate
// incoming requests, inject user information into the context, and enforce security policies
// before passing control to the actual gRPC handlers.
package interceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthUnaryInterceptor is a gRPC unary server interceptor that handles authentication via metadata tokens.
// It checks for the presence of an "authorization" token in the incoming context metadata. If the token is valid,
// it extracts the associated user information and appends it to the metadata. If the token is missing or invalid,
// it generates a new user ID and token, appends them to the metadata, and updates the context accordingly.
// The interceptor then calls the handler with the updated context. Returns an error if metadata is missing or
// token creation fails.
func AuthUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata is not provided")
	}

	tokenFromMD := md.Get("authorization")

	if len(tokenFromMD) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	ok, user := auth.Validate(tokenFromMD[0])

	var newToken string

	if !ok {
		generatedUserID := uuid.New()

		var err error
		newToken, err = auth.CreateToken(generatedUserID)

		if err != nil {
			return nil, status.Errorf(codes.Internal, "error creating new token: %v", err)
		}
		md.Append("user_id", generatedUserID.String())
	}

	md.Append("user_id", user.String())
	md.Append("token", newToken)

	ctx = metadata.NewIncomingContext(ctx, md)

	return handler(ctx, req)
}

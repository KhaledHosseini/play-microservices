package interceptors

import (
	"context"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"

	"github.com/golang-jwt/jwt/v5"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	log logger.Logger
	cfg *config.Config
}

func NewAuthInterceptor(log logger.Logger, cfg *config.Config) *AuthInterceptor {
	return &AuthInterceptor{log: log, cfg: cfg}
}

func (ai *AuthInterceptor) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract the Authorization header from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}
	authHeaders, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Missing Authorization header")
	}

	tokenString := authHeaders[0]
	// Verify and parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != jwt.SigningMethodRS256.Name {
			return nil, status.Error(codes.Unauthenticated, "Invalid algorithm")
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(ai.cfg.AuthPublicKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
		ai.log.Info("Public key is: %v", verifyKey)
		return verifyKey, nil
	})

	if err != nil {
		ai.log.Errorf("Failed to parse JWT token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ai.log.Info(claims["nbf"])
		// Access the role value
		role := claims["role"]
		ai.log.Info(role)
		//  We can Add the role information to the metadata of the request to be used
		// by the service method for authorization puroposes. but we do a simple authorization here!
		if role != "admin" {
			return nil, status.Errorf(codes.PermissionDenied, "Only admins have access!")
		}
		// newMD := metadata.New(map[string]string{"role": role})
		// ctx = metadata.NewIncomingContext(ctx, newMD)
		// Call the actual gRPC handler with updated context
		return handler(ctx, req)
	} else {
		ai.log.Info(err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

}

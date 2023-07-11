package interceptors

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// We do authentication and authorization in the end services. We just attach the auth headers to grpc requests.
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization") // Retrieve the access token from the request header

		// do authentication.

		// Create a gRPC context and add the access token to the metadata
		md := metadata.Pairs("authorization", accessToken)
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		// Update the request context with the modified context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

package interceptors

import (
	"context"
	"net/http"

	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/pkg/cookie"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// We do authentication and authorization in the end services. We just attach the auth headers to grpc requests.
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := cookie.GetAccessToken(c) // Retrieve the access token from the request header
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": err.Error()})
			return
		}
		// do authentication if necessary.
		// In Microservices architecture,
		// One approach is to do authentication inside api-gateway only and do the authorization in downstream services.
		// Another aproach is to do both authentication and authorization in each service.

		// gRPC methods can’t set HTTP response cookies directly —
		// they’re running in a different process and don’t have access to the response object.
		// Create a gRPC context and add the access token to the metadata see: https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
		md := metadata.New(map[string]string{"authorization": accessToken})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

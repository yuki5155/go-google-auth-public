package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	ginpkg "github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/handlers"
	"github.com/yuki5155/go-google-auth/internal/services"
)

var ginLambda *ginadapter.GinLambda

func init() {
	cfg := config.Load()
	jwtService := services.NewJWTService(cfg.JWTSecret)

	if cfg.IsProduction() {
		ginpkg.SetMode(ginpkg.ReleaseMode)
	}

	r := ginpkg.Default()
	r.Use(corsMiddleware(cfg))

	authHandler := handlers.NewAuthHandler(cfg, jwtService)
	r.POST("/auth/refresh", authHandler.RefreshToken)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

func corsMiddleware(cfg *config.Config) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range cfg.AllowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if len(cfg.AllowedOrigins) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.AllowedOrigins[0])
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

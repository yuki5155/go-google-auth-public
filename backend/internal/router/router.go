package router

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/handlers"
	"github.com/yuki5155/go-google-auth/internal/middleware"
	"github.com/yuki5155/go-google-auth/internal/services"
)

// Setup initializes and configures the Gin router with all routes and middleware
func Setup(cfg *config.Config, jwtService *services.JWTService) *gin.Engine {
	// 本番環境ではGinのリリースモードを使用
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Ginルーターの初期化
	r := gin.Default()

	// CORS設定（環境変数から）
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 許可されたOriginかチェック
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
			// デフォルトで最初のOriginを使用
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
	})

	// ハンドラーの初期化
	helloHandler := handlers.NewHelloHandler()
	healthHandler := handlers.NewHealthHandler()
	setCookieHandler := handlers.NewCookieHandler(cfg)
	checkCookieHandler := handlers.NewCheckCookieHandler()
	authHandler := handlers.NewAuthHandler(cfg, jwtService)

	// 公開ルーティング設定
	r.GET(helloHandler.Path, helloHandler.Handle)
	r.GET(healthHandler.Path, healthHandler.Handle)
	r.GET("/health/ready", healthHandler.Handle)

	// Cookie test endpoints
	r.GET(setCookieHandler.Path, setCookieHandler.Handle)
	r.GET(checkCookieHandler.Path, checkCookieHandler.Handle)

	// Auth endpoints (public)
	r.POST("/auth/google", authHandler.GoogleLogin)
	r.POST("/auth/refresh", authHandler.RefreshToken)
	r.POST("/auth/logout", authHandler.Logout)

	// Protected routes (require authentication)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		protected.GET("/me", authHandler.GetCurrentUser)
	}

	log.Printf("Router configured (environment: %s)", cfg.Environment)
	log.Printf("Allowed CORS origins: %s", strings.Join(cfg.AllowedOrigins, ", "))
	log.Printf("Google Client ID configured: %v", cfg.GoogleClientID != "")

	return r
}

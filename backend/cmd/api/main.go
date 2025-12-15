package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/handlers"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

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

	// ハンドラーの初期化（configを渡す）
	helloHandler := handlers.NewHelloHandler()
	healthHandler := handlers.NewHealthHandler()
	setCookieHandler := handlers.NewCookieHandler(cfg)
	checkCookieHandler := handlers.NewCheckCookieHandler()

	// ルーティング設定
	r.GET(helloHandler.Path, helloHandler.Handle)
	r.GET(healthHandler.Path, healthHandler.Handle)
	r.GET("/health/ready", healthHandler.Handle)

	// Cookie test endpoints
	r.GET(setCookieHandler.Path, setCookieHandler.Handle)
	r.GET(checkCookieHandler.Path, checkCookieHandler.Handle)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Environment)
	log.Printf("Allowed CORS origins: %s", strings.Join(cfg.AllowedOrigins, ", "))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

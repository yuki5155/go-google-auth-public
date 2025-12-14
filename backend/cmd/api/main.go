package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/handlers"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

	// Ginルーターの初期化
	r := gin.Default()

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
	setCookieHandler := handlers.NewCookieHandler()
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
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

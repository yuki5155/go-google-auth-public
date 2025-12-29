package main

import (
	"fmt"
	"log"

	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/router"
	"github.com/yuki5155/go-google-auth/internal/services"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

	// サービスの初期化
	jwtService := services.NewJWTService(cfg.JWTSecret)

	// ルーターのセットアップ
	r := router.Setup(cfg, jwtService)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

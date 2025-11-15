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

	// ハンドラーの初期化
	helloHandler := handlers.NewHelloHandler()

	// ルーティング設定
	r.GET(helloHandler.Path, helloHandler.Handle)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

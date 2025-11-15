.PHONY: up down build logs restart clean help

# Docker Composeを起動
up:
	docker compose up -d

# Docker Composeを停止
down:
	docker compose down

# Docker Composeをビルドして起動
build:
	docker compose up -d --build

# ログを表示
logs:
	docker compose logs -f app

# 再起動
restart:
	docker compose restart app

# コンテナとボリュームを完全削除
clean:
	docker compose down -v

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo "  make up       - Docker Composeを起動"
	@echo "  make down     - Docker Composeを停止"
	@echo "  make build    - Docker Composeをビルドして起動"
	@echo "  make logs     - アプリケーションのログを表示"
	@echo "  make restart  - アプリケーションを再起動"
	@echo "  make clean    - コンテナとボリュームを完全削除"
	@echo "  make help     - このヘルプを表示"


.PHONY: up down build logs restart clean env help

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

# .envファイルを生成（既存の場合はスキップ）- Docker経由で実行
env:
	@docker run --rm -v $(PWD):/app -w /app alpine sh -c '\
		if [ ! -f .env ]; then \
			cp .env.example .env; \
			echo ".env file created from .env.example"; \
			echo "Please update the values in .env with your actual credentials"; \
		else \
			echo ".env file already exists. Remove it first if you want to regenerate."; \
		fi'

# .envファイルを強制的に再生成 - Docker経由で実行
env-force:
	@docker run --rm -v $(PWD):/app -w /app alpine sh -c 'cp .env.example .env && echo ".env file created from .env.example"'

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo "  make up        - Docker Composeを起動"
	@echo "  make down      - Docker Composeを停止"
	@echo "  make build     - Docker Composeをビルドして起動"
	@echo "  make logs      - アプリケーションのログを表示"
	@echo "  make restart   - アプリケーションを再起動"
	@echo "  make clean     - コンテナとボリュームを完全削除"
	@echo "  make env       - .envファイルを生成（Docker経由、既存の場合はスキップ）"
	@echo "  make env-force - .envファイルを強制的に再生成（Docker経由）"
	@echo "  make help      - このヘルプを表示"


# Makefile

# Variables
DC = Dockerfile
GO = go
IMAGE_NAME = image-gallery-app

.PHONY: help
help:
	@echo "Makefile for building and running a Go application in Docker"
	@echo "Usage:"
	@echo "  make all            : 実行ファイルの作成"
	@echo "  make go-build       : 実行ファイルの作成"
	@echo "  make run            : ローカルで起動"
	@echo "  make help           : ヘルプを表示"
	@echo "  make build-image    : Dockerイメージのビルド"
	@echo "  make down           : Dockerコンテナの停止"
	@echo "  make up             : Dockerコンテナの起動"
	@echo "  make clean          : Dockerコンテナのクリーンアップ"
	@echo "  make migrate        : DBマイグレーションの実行"
	@echo "  make test           : テストの実行"
	@echo "  make lint           : Lintチェックの実行"

.PHONY: all
all: go-build

.PHONY: go-build
all: go-build
go-build:
	@echo "=== 実行ファイルの作成 ==="
	$(GO) build -o ./backend/app ./backend/cmd
	@echo "Docker image built successfully."

.PHONY: run
run:
	@echo "=== ローカルで起動 ==="
	$(GO) run ./backend/cmd
	@echo "Docker container started successfully."

.PHONY: build-image
build-image:
	@echo "=== Dockerイメージのビルド: $(IMAGE_NAME) ==="
	docker-compose build

.PHONY: down
down:
	@echo "=== Dockerコンテナの停止 ==="
	@docker-compose down
	@echo "Docker container stopped successfully."

.PHONY: up
up:
	@echo "=== Dockerコンテナの起動 ==="
	@docker-compose up -d
	@echo "Docker container started successfully."

.PHONY: clean
clean:
	@echo "=== Dockerコンテナのクリーンアップ ==="
	@docker-compose down
	@docker rmi my-go-app
	@rm -rf ./backend/app
	@echo "Cleaned up successfully."

.PHONY: migrate
migrate:
	@echo "=== DBマイグレーションの実行 ==="
	@docker-compose exec -T db sh -c "cd /app && go run ./migrations/migrate.go"
	@echo "Database migration executed successfully."

.PHONY: test
test:
	@echo "=== テストの実行 ==="
	@docker-compose exec -T app sh -c "cd /app && go test ./..."
	@echo "Tests executed successfully."

.PHONY: lint
lint:
	@echo "=== Lintチェックの実行 ==="
	@docker-compose exec -T app sh -c "cd /app && golangci-lint run"
	@echo "Lint check executed successfully."

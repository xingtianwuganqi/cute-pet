# 项目名称
APP_NAME = pet-project
# 构建输出目录
BUILD_DIR = ./bin
# 二进制文件名
BINARY = $(BUILD_DIR)/$(APP_NAME)
# Docker镜像名称
DOCKER_IMAGE = cute-pet/pet-project
# Docker镜像标签
DOCKER_TAG = latest
# Docker registry地址（可选）
DOCKER_REGISTRY ?= 

# 获取当前系统信息
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

# 默认目标
.PHONY: all
all: build

# 构建项目
.PHONY: build
build:
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) .

# 运行项目
.PHONY: run
run: build
	@echo "Running $(APP_NAME)"
	./$(BINARY)

# 开发模式运行（带热重载）
.PHONY: dev
dev:
	@echo "Starting development mode..."
	@go run main.go

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 检查代码错误
.PHONY: vet
vet:
	@echo "Checking code with go vet..."
	@go vet ./...

# Docker构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
ifeq ($(DOCKER_REGISTRY), )
	@echo "Docker image $(DOCKER_IMAGE):$(DOCKER_TAG) built successfully."
else
	@echo "Docker image $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) built successfully."
endif

# Docker推送
.PHONY: docker-push
docker-push: docker-build
	@echo "Tagging image for registry..."
ifeq ($(DOCKER_REGISTRY), )
	@echo "DOCKER_REGISTRY is not set, skipping push. Please set DOCKER_REGISTRY environment variable."
	@exit 1
else
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Pushing image to $(DOCKER_REGISTRY)..."
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Image pushed to $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)"
endif

# Docker运行
.PHONY: docker-run
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -d \
		--name $(APP_NAME) \
		-p 8080:8080 \
		-e ENVIRONMENT=production \
		-v $(PWD)/logs:/app/logs \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker停止
.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true

# Docker清理
.PHONY: docker-clean
docker-clean:
	@echo "Removing Docker container and image..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true

# Docker重新构建并运行
.PHONY: docker-restart
docker-restart: docker-stop docker-build docker-run
	@echo "Docker container restarted."

# Docker查看日志
.PHONY: docker-logs
docker-logs:
	@docker logs -f $(APP_NAME)

# Docker进入容器
.PHONY: docker-exec
docker-exec:
	@docker exec -it $(APP_NAME) /bin/sh

# Docker查看状态
.PHONY: docker-ps
docker-ps:
	@docker ps | grep $(APP_NAME)

# 针对远程服务器的构建
.PHONY: remote-build
remote-build:
	@echo "Building for Linux platform..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARY)-linux-amd64 .

# 针对远程服务器的Docker构建
.PHONY: remote-docker-build
remote-docker-build:
	@echo "Building Docker image for remote deployment..."
	@docker buildx build --platform linux/amd64 -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# 多平台构建
.PHONY: multi-platform-build
multi-platform-build:
	@echo "Building multi-platform Docker image..."
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):$(DOCKER_TAG) --push .

# 部署到远程服务器
.PHONY: deploy
deploy: docker-build
	@echo "Deploying to remote server..."
ifeq ($(REMOTE_HOST), )
	@echo "REMOTE_HOST is not set. Usage: make deploy REMOTE_HOST=user@hostname"
	@exit 1
else
	@echo "Deploying to $(REMOTE_HOST)"
	scp Dockerfile .dockerignore *.go go.* config/ $(REMOTE_HOST):~/$(APP_NAME)/
	ssh $(REMOTE_HOST) "cd ~/$(APP_NAME) && docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) . && docker run -d --name $(APP_NAME) -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)"
endif

# 环境变量设置
.PHONY: set-env
set-env:
	@echo "Setting environment variables..."
	@export ENVIRONMENT=production

# 显示帮助信息
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Build the project (default)"
	@echo "  build            Build the binary"
	@echo "  run              Run the built binary"
	@echo "  dev              Run in development mode"
	@echo "  test             Run tests"
	@echo "  deps             Install dependencies"
	@echo "  clean            Clean build artifacts"
	@echo "  fmt              Format code"
	@echo "  vet              Vet code"
	@echo "  docker-build     Build Docker image"
	@echo "  docker-push      Push Docker image to registry"
	@echo "  docker-run       Run in Docker container"
	@echo "  docker-stop      Stop Docker container"
	@echo "  docker-clean     Remove Docker container and image"
	@echo "  docker-restart   Restart Docker container"
	@echo "  docker-logs      View Docker container logs"
	@echo "  docker-exec      Execute command in Docker container"
	@echo "  docker-ps        List running Docker containers"
	@echo "  remote-build     Build for remote Linux server"
	@echo "  remote-docker-build Build Docker image for remote deployment"
	@echo "  multi-platform-build Build multi-platform Docker image"
	@echo "  deploy           Deploy to remote server (requires REMOTE_HOST)"
	@echo "  set-env          Set environment variables"
	@echo "  help             Show this help message"
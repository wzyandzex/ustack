# ustack Makefile

.PHONY: all build clean test lint client server

# 默认目标
all: build

# 构建所有目标
build: client server

# 构建客户端
client:
	@echo "Building ustack-client..."
	@mkdir -p bin
	go build -o bin/ustack-client ./cmd/client

# 构建服务端
server:
	@echo "Building ustack-server..."
	@mkdir -p bin
	go build -o bin/ustack-server ./cmd/server

# 运行测试
test:
	@echo "Running tests..."
	go test ./... -v

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 运行服务端示例
run-server: server
	@echo "Running ustack-server on port 8080..."
	./bin/ustack-server 8080

# 运行客户端示例
run-client: client
	@echo "Running ustack-client..."
	./bin/ustack-client localhost 8080

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build          - Build client and server"
	@echo "  client         - Build client only"
	@echo "  server         - Build server only"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  clean          - Clean build files"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  run-server     - Run server example"
	@echo "  run-client     - Run client example"
	@echo "  help           - Show this help" 
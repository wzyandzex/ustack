package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"ustack/internal/utils"
	"ustack/pkg/tcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ustack-server <port>")
		fmt.Println("Example: ustack-server 8080")
		os.Exit(1)
	}

	portStr := os.Args[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("Invalid port number: %s\n", portStr)
		os.Exit(1)
	}

	logger := utils.DefaultLogger
	logger.Info("Starting ustack HTTP server on port %d...", port)

	// 创建本地IP
	var localIP [4]byte
	copy(localIP[:], net.ParseIP("127.0.0.1").To4())

	// 创建TCP连接（监听模式）
	conn := tcp.NewConnection(localIP, uint16(port), [4]byte{}, 0)

	// 设置回调函数
	conn.OnStateChanged = func(state string) {
		logger.Info("Server state changed: %s", state)
	}

	conn.OnDataReceived = func(data []byte) {
		logger.Info("Received HTTP request: %d bytes", len(data))

		// 解析HTTP请求
		request := string(data)
		logger.Info("HTTP Request:\n%s", request)

		// 生成HTTP响应
		response := generateHTTPResponse()

		// 发送响应
		err := conn.Send([]byte(response))
		if err != nil {
			logger.Error("Failed to send HTTP response: %v", err)
		} else {
			logger.Info("HTTP response sent: %d bytes", len(response))
		}
	}

	// 开始监听
	err = conn.Listen()
	if err != nil {
		logger.Error("Failed to start listening: %v", err)
		os.Exit(1)
	}

	logger.Info("Server is listening on port %d", port)
	logger.Info("Press Ctrl+C to stop the server")

	// 保持服务器运行（简化处理）
	select {}
}

// generateHTTPResponse 生成HTTP响应
func generateHTTPResponse() string {
	return `HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 89
Connection: close

<!DOCTYPE html>
<html>
<head><title>ustack HTTP Server</title></head>
<body>
<h1>Hello from ustack!</h1>
<p>This is a response from the user-space TCP/IP stack.</p>
</body>
</html>`
}

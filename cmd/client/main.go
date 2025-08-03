package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"ustack/internal/utils"
	"ustack/pkg/tcp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ustack-client <host> <port>")
		fmt.Println("Example: ustack-client localhost 8080")
		os.Exit(1)
	}

	host := os.Args[1]
	port := os.Args[2]

	logger := utils.DefaultLogger
	logger.Info("Starting ustack HTTP client...")
	logger.Info("Connecting to %s:%s", host, port)

	// 解析目标地址
	remoteIP := net.ParseIP(host)
	if remoteIP == nil {
		logger.Error("Invalid host address: %s", host)
		os.Exit(1)
	}

	var remoteIPBytes [4]byte
	copy(remoteIPBytes[:], remoteIP.To4())

	// 创建本地IP（简化处理）
	var localIP [4]byte
	copy(localIP[:], net.ParseIP("127.0.0.1").To4())

	// 创建TCP连接
	conn := tcp.NewConnection(localIP, 12345, remoteIPBytes, 8080)

	// 设置回调函数
	conn.OnStateChanged = func(state string) {
		logger.Info("Connection state changed: %s", state)
	}

	conn.OnDataReceived = func(data []byte) {
		logger.Info("Received data: %d bytes", len(data))
		fmt.Printf("Response:\n%s\n", string(data))
	}

	// 建立连接
	err := conn.Connect()
	if err != nil {
		logger.Error("Failed to connect: %v", err)
		os.Exit(1)
	}

	logger.Info("Connection established successfully")

	// 发送HTTP GET请求
	httpRequest := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s:%s\r\nConnection: close\r\n\r\n", host, port)

	err = conn.Send([]byte(httpRequest))
	if err != nil {
		logger.Error("Failed to send HTTP request: %v", err)
		os.Exit(1)
	}

	logger.Info("HTTP request sent")

	// 等待响应（简化处理）
	fmt.Println("Press Enter to close connection...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	// 关闭连接
	err = conn.Close()
	if err != nil {
		logger.Error("Failed to close connection: %v", err)
	}

	logger.Info("Connection closed")
}

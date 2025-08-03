package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
	"ustack/internal/utils"
)

const (
	// TCP状态
	StateClosed      = "CLOSED"
	StateListen      = "LISTEN"
	StateSynSent     = "SYN_SENT"
	StateSynReceived = "SYN_RECEIVED"
	StateEstablished = "ESTABLISHED"
	StateFinWait1    = "FIN_WAIT_1"
	StateFinWait2    = "FIN_WAIT_2"
	StateCloseWait   = "CLOSE_WAIT"
	StateClosing     = "CLOSING"
	StateLastAck     = "LAST_ACK"
	StateTimeWait    = "TIME_WAIT"

	// 默认窗口大小
	DefaultWindowSize = 65535

	// 超时时间
	ConnectionTimeout = 30 * time.Second
	RetransmitTimeout = 3 * time.Second
)

// Connection TCP连接结构
type Connection struct {
	mu sync.Mutex

	// 连接标识
	LocalIP    [4]byte
	LocalPort  uint16
	RemoteIP   [4]byte
	RemotePort uint16

	// 状态
	State string

	// 序列号
	SendSequence    uint32
	ReceiveSequence uint32

	// 窗口
	SendWindow    uint16
	ReceiveWindow uint16

	// 缓冲区
	SendBuffer    []byte
	ReceiveBuffer []byte

	// 拥塞控制
	CongestionWindow   uint16
	SlowStartThreshold uint16

	// 定时器
	RetransmitTimer *time.Timer
	KeepAliveTimer  *time.Timer

	// 回调函数
	OnDataReceived func([]byte)
	OnStateChanged func(string)

	// 日志
	logger *utils.Logger
}

// NewConnection 创建新的TCP连接
func NewConnection(localIP [4]byte, localPort uint16, remoteIP [4]byte, remotePort uint16) *Connection {
	conn := &Connection{
		LocalIP:            localIP,
		LocalPort:          localPort,
		RemoteIP:           remoteIP,
		RemotePort:         remotePort,
		State:              StateClosed,
		SendWindow:         DefaultWindowSize,
		ReceiveWindow:      DefaultWindowSize,
		CongestionWindow:   1, // 慢启动初始窗口
		SlowStartThreshold: 65535,
		SendBuffer:         make([]byte, 0, 8192),
		ReceiveBuffer:      make([]byte, 0, 8192),
		logger:             utils.DefaultLogger,
	}

	return conn
}

// Connect 建立连接（客户端）
func (c *Connection) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.State != StateClosed {
		return fmt.Errorf("connection not in CLOSED state")
	}

	// 生成随机序列号
	c.SendSequence = rand.Uint32()
	c.ReceiveSequence = 0

	// 发送SYN包
	synHeader := NewHeader(c.LocalPort, c.RemotePort, c.SendSequence, 0, FlagSYN, c.SendWindow)

	c.State = StateSynSent
	c.logger.LogConnection("SYN_SENT", fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort),
		fmt.Sprintf("%s:%d", net.IP(c.RemoteIP[:]), c.RemotePort))

	// 这里应该发送SYN包到网络层
	// 简化处理，直接模拟收到SYN+ACK
	c.handleSynAck(synHeader)

	return nil
}

// Listen 监听连接（服务端）
func (c *Connection) Listen() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.State = StateListen
	c.logger.LogConnection("LISTEN", fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort), "")

	return nil
}

// Send 发送数据
func (c *Connection) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.State != StateEstablished {
		return fmt.Errorf("connection not established")
	}

	// 添加到发送缓冲区
	c.SendBuffer = append(c.SendBuffer, data...)

	// 创建数据包（这里应该通过IP层发送）
	_ = NewHeader(c.LocalPort, c.RemotePort, c.SendSequence, c.ReceiveSequence, FlagPSH|FlagACK, c.SendWindow)

	// 发送数据
	// 这里应该通过IP层发送
	c.SendSequence += uint32(len(data))

	c.logger.LogPacket("SEND", "TCP", fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort),
		fmt.Sprintf("%s:%d", net.IP(c.RemoteIP[:]), c.RemotePort), len(data))

	return nil
}

// Receive 接收数据
func (c *Connection) Receive(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 添加到接收缓冲区
	c.ReceiveBuffer = append(c.ReceiveBuffer, data...)

	// 更新接收序列号
	c.ReceiveSequence += uint32(len(data))

	// 发送ACK（这里应该通过IP层发送）
	_ = NewHeader(c.LocalPort, c.RemotePort, c.SendSequence, c.ReceiveSequence, FlagACK, c.ReceiveWindow)

	c.logger.LogPacket("RECV", "TCP", fmt.Sprintf("%s:%d", net.IP(c.RemoteIP[:]), c.RemotePort),
		fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort), len(data))

	// 调用数据接收回调
	if c.OnDataReceived != nil {
		c.OnDataReceived(data)
	}

	return nil
}

// Close 关闭连接
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.State == StateClosed {
		return nil
	}

	// 发送FIN包（这里应该通过IP层发送）
	_ = NewHeader(c.LocalPort, c.RemotePort, c.SendSequence, c.ReceiveSequence, FlagFIN|FlagACK, c.SendWindow)

	c.State = StateFinWait1
	c.logger.LogConnection("FIN_WAIT_1", fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort),
		fmt.Sprintf("%s:%d", net.IP(c.RemoteIP[:]), c.RemotePort))

	// 简化处理，直接关闭
	c.State = StateClosed

	return nil
}

// handleSynAck 处理SYN+ACK包
func (c *Connection) handleSynAck(synHeader *Header) {
	// 模拟收到SYN+ACK
	c.ReceiveSequence = synHeader.SequenceNumber + 1

	// 发送ACK（这里应该通过IP层发送）
	_ = NewHeader(c.LocalPort, c.RemotePort, c.SendSequence, c.ReceiveSequence, FlagACK, c.SendWindow)

	c.State = StateEstablished
	c.logger.LogConnection("ESTABLISHED", fmt.Sprintf("%s:%d", net.IP(c.LocalIP[:]), c.LocalPort),
		fmt.Sprintf("%s:%d", net.IP(c.RemoteIP[:]), c.RemotePort))

	if c.OnStateChanged != nil {
		c.OnStateChanged(c.State)
	}
}

// String 返回连接的字符串表示
func (c *Connection) String() string {
	return fmt.Sprintf("TCP Connection: %s:%d -> %s:%d [%s]",
		net.IP(c.LocalIP[:]), c.LocalPort,
		net.IP(c.RemoteIP[:]), c.RemotePort,
		c.State)
}

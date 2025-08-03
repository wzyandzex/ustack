package icmp

import (
	"encoding/binary"
	"fmt"
	"ustack/internal/utils"
)

const (
	// ICMP类型
	TypeEchoRequest  = 8
	TypeEchoReply    = 0
	TypeDestUnreach  = 3
	TypeTimeExceeded = 11

	// ICMP代码
	CodeEchoRequest = 0
	CodeEchoReply   = 0
)

// Packet ICMP数据包结构
type Packet struct {
	Type     uint8  // 类型
	Code     uint8  // 代码
	Checksum uint16 // 校验和
	ID       uint16 // 标识符
	Sequence uint16 // 序列号
	Data     []byte // 数据
}

// Marshal 将ICMP数据包序列化为字节数组
func (p *Packet) Marshal() ([]byte, error) {
	// 计算数据长度（头部8字节 + 数据）
	dataLength := 8 + len(p.Data)
	data := make([]byte, dataLength)

	// 类型
	data[0] = p.Type

	// 代码
	data[1] = p.Code

	// 校验和（先设为0）
	binary.BigEndian.PutUint16(data[2:4], 0)

	// 标识符
	binary.BigEndian.PutUint16(data[4:6], p.ID)

	// 序列号
	binary.BigEndian.PutUint16(data[6:8], p.Sequence)

	// 数据
	copy(data[8:], p.Data)

	// 计算校验和
	p.Checksum = utils.CalculateChecksum(data)
	binary.BigEndian.PutUint16(data[2:4], p.Checksum)

	return data, nil
}

// Unmarshal 从字节数组解析ICMP数据包
func (p *Packet) Unmarshal(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("ICMP packet too short: %d bytes", len(data))
	}

	// 类型
	p.Type = data[0]

	// 代码
	p.Code = data[1]

	// 校验和
	p.Checksum = binary.BigEndian.Uint16(data[2:4])

	// 标识符
	p.ID = binary.BigEndian.Uint16(data[4:6])

	// 序列号
	p.Sequence = binary.BigEndian.Uint16(data[6:8])

	// 数据
	p.Data = make([]byte, len(data)-8)
	copy(p.Data, data[8:])

	return nil
}

// String 返回ICMP数据包的字符串表示
func (p *Packet) String() string {
	return fmt.Sprintf("ICMP Packet: Type=%d, Code=%d, ID=%d, Sequence=%d, Data=%d bytes",
		p.Type, p.Code, p.ID, p.Sequence, len(p.Data))
}

// IsEchoRequest 检查是否为Echo Request
func (p *Packet) IsEchoRequest() bool {
	return p.Type == TypeEchoRequest
}

// IsEchoReply 检查是否为Echo Reply
func (p *Packet) IsEchoReply() bool {
	return p.Type == TypeEchoReply
}

// CreateReply 创建回复数据包
func (p *Packet) CreateReply() *Packet {
	return &Packet{
		Type:     TypeEchoReply,
		Code:     CodeEchoReply,
		Checksum: 0, // 将由Marshal计算
		ID:       p.ID,
		Sequence: p.Sequence,
		Data:     p.Data, // 复制原始数据
	}
}

// NewEchoRequest 创建新的Echo Request数据包
func NewEchoRequest(id, sequence uint16, data []byte) *Packet {
	return &Packet{
		Type:     TypeEchoRequest,
		Code:     CodeEchoRequest,
		Checksum: 0, // 将由Marshal计算
		ID:       id,
		Sequence: sequence,
		Data:     data,
	}
}

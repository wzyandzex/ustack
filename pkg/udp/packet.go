package udp

import (
	"encoding/binary"
	"fmt"
	"ustack/internal/utils"
)

const (
	// UDP头部长度
	UDPHeaderLength = 8
)

// Packet UDP数据包结构
type Packet struct {
	SourcePort      uint16 // 源端口
	DestinationPort uint16 // 目标端口
	Length          uint16 // 长度
	Checksum        uint16 // 校验和
	Payload         []byte // 数据载荷
}

// Marshal 将UDP数据包序列化为字节数组
func (p *Packet) Marshal() ([]byte, error) {
	// 计算总长度（头部8字节 + 数据）
	totalLength := UDPHeaderLength + len(p.Payload)
	if totalLength > 65535 {
		return nil, fmt.Errorf("UDP packet too large: %d bytes", totalLength)
	}

	data := make([]byte, totalLength)

	// 源端口
	binary.BigEndian.PutUint16(data[0:2], p.SourcePort)

	// 目标端口
	binary.BigEndian.PutUint16(data[2:4], p.DestinationPort)

	// 长度
	binary.BigEndian.PutUint16(data[4:6], uint16(totalLength))

	// 校验和（先设为0）
	binary.BigEndian.PutUint16(data[6:8], 0)

	// 数据载荷
	copy(data[8:], p.Payload)

	// 计算校验和（这里简化处理，实际应该包含伪头部）
	p.Checksum = utils.CalculateChecksum(data)
	binary.BigEndian.PutUint16(data[6:8], p.Checksum)

	return data, nil
}

// Unmarshal 从字节数组解析UDP数据包
func (p *Packet) Unmarshal(data []byte) error {
	if len(data) < UDPHeaderLength {
		return fmt.Errorf("UDP packet too short: %d bytes", len(data))
	}

	// 源端口
	p.SourcePort = binary.BigEndian.Uint16(data[0:2])

	// 目标端口
	p.DestinationPort = binary.BigEndian.Uint16(data[2:4])

	// 长度
	p.Length = binary.BigEndian.Uint16(data[4:6])

	// 校验和
	p.Checksum = binary.BigEndian.Uint16(data[6:8])

	// 数据载荷
	payloadLength := len(data) - UDPHeaderLength
	p.Payload = make([]byte, payloadLength)
	copy(p.Payload, data[8:])

	return nil
}

// String 返回UDP数据包的字符串表示
func (p *Packet) String() string {
	return fmt.Sprintf("UDP Packet: %d -> %d, Length: %d, Payload: %d bytes",
		p.SourcePort, p.DestinationPort, p.Length, len(p.Payload))
}

// NewPacket 创建新的UDP数据包
func NewPacket(srcPort, dstPort uint16, payload []byte) *Packet {
	return &Packet{
		SourcePort:      srcPort,
		DestinationPort: dstPort,
		Length:          uint16(UDPHeaderLength + len(payload)),
		Checksum:        0, // 将由Marshal计算
		Payload:         payload,
	}
}

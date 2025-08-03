package eth

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// 以太网帧头部长度
	EthernetHeaderLength = 14

	// 以太网类型
	EtherTypeIPv4 = 0x0800
	EtherTypeARP  = 0x0806
)

// Frame 以太网帧结构
type Frame struct {
	DestinationMAC [6]byte // 目标MAC地址
	SourceMAC      [6]byte // 源MAC地址
	EtherType      uint16  // 以太网类型
	Payload        []byte  // 数据载荷
}

// Marshal 将帧序列化为字节数组
func (f *Frame) Marshal() ([]byte, error) {
	if len(f.Payload) > 1500 {
		return nil, fmt.Errorf("payload too large: %d bytes", len(f.Payload))
	}

	data := make([]byte, EthernetHeaderLength+len(f.Payload))

	// 复制目标MAC地址
	copy(data[0:6], f.DestinationMAC[:])

	// 复制源MAC地址
	copy(data[6:12], f.SourceMAC[:])

	// 设置以太网类型
	binary.BigEndian.PutUint16(data[12:14], f.EtherType)

	// 复制载荷
	copy(data[14:], f.Payload)

	return data, nil
}

// Unmarshal 从字节数组解析帧
func (f *Frame) Unmarshal(data []byte) error {
	if len(data) < EthernetHeaderLength {
		return fmt.Errorf("frame too short: %d bytes", len(data))
	}

	// 解析目标MAC地址
	copy(f.DestinationMAC[:], data[0:6])

	// 解析源MAC地址
	copy(f.SourceMAC[:], data[6:12])

	// 解析以太网类型
	f.EtherType = binary.BigEndian.Uint16(data[12:14])

	// 解析载荷
	f.Payload = make([]byte, len(data)-EthernetHeaderLength)
	copy(f.Payload, data[14:])

	return nil
}

// String 返回帧的字符串表示
func (f *Frame) String() string {
	return fmt.Sprintf("Ethernet Frame: %s -> %s, Type: 0x%04x, Payload: %d bytes",
		net.HardwareAddr(f.SourceMAC[:]).String(),
		net.HardwareAddr(f.DestinationMAC[:]).String(),
		f.EtherType,
		len(f.Payload))
}

// IsBroadcast 检查是否为广播帧
func (f *Frame) IsBroadcast() bool {
	return f.DestinationMAC == [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
}

// IsMulticast 检查是否为多播帧
func (f *Frame) IsMulticast() bool {
	return (f.DestinationMAC[0] & 0x01) != 0
}

// NewFrame 创建新的以太网帧
func NewFrame(srcMAC, dstMAC [6]byte, etherType uint16, payload []byte) *Frame {
	return &Frame{
		SourceMAC:      srcMAC,
		DestinationMAC: dstMAC,
		EtherType:      etherType,
		Payload:        payload,
	}
}

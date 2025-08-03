package tcp

import (
	"encoding/binary"
	"fmt"
	"ustack/internal/utils"
)

const (
	// TCP头部长度
	TCPHeaderLength = 20

	// TCP标志
	FlagFIN = 0x01
	FlagSYN = 0x02
	FlagRST = 0x04
	FlagPSH = 0x08
	FlagACK = 0x10
	FlagURG = 0x20
)

// Header TCP头部结构
type Header struct {
	SourcePort      uint16 // 源端口
	DestinationPort uint16 // 目标端口
	SequenceNumber  uint32 // 序列号
	Acknowledgment  uint32 // 确认号
	DataOffset      uint8  // 数据偏移
	Flags           uint8  // 标志
	WindowSize      uint16 // 窗口大小
	Checksum        uint16 // 校验和
	UrgentPointer   uint16 // 紧急指针
	Options         []byte // 选项（可选）
}

// Marshal 将TCP头部序列化为字节数组
func (h *Header) Marshal() ([]byte, error) {
	// 计算头部长度（基本头部20字节 + 选项）
	headerLength := TCPHeaderLength + len(h.Options)
	if headerLength > 60 {
		return nil, fmt.Errorf("TCP header too large: %d bytes", headerLength)
	}

	data := make([]byte, headerLength)

	// 源端口
	binary.BigEndian.PutUint16(data[0:2], h.SourcePort)

	// 目标端口
	binary.BigEndian.PutUint16(data[2:4], h.DestinationPort)

	// 序列号
	binary.BigEndian.PutUint32(data[4:8], h.SequenceNumber)

	// 确认号
	binary.BigEndian.PutUint32(data[8:12], h.Acknowledgment)

	// 数据偏移和标志
	dataOffset := uint8(headerLength / 4) // 以4字节为单位
	data[12] = (dataOffset << 4) | (h.Flags & 0x3F)

	// 窗口大小
	binary.BigEndian.PutUint16(data[14:16], h.WindowSize)

	// 校验和（先设为0）
	binary.BigEndian.PutUint16(data[16:18], 0)

	// 紧急指针
	binary.BigEndian.PutUint16(data[18:20], h.UrgentPointer)

	// 选项
	if len(h.Options) > 0 {
		copy(data[20:], h.Options)
	}

	// 计算校验和（这里简化处理，实际应该包含伪头部）
	h.Checksum = utils.CalculateChecksum(data)
	binary.BigEndian.PutUint16(data[16:18], h.Checksum)

	return data, nil
}

// Unmarshal 从字节数组解析TCP头部
func (h *Header) Unmarshal(data []byte) error {
	if len(data) < TCPHeaderLength {
		return fmt.Errorf("TCP header too short: %d bytes", len(data))
	}

	// 源端口
	h.SourcePort = binary.BigEndian.Uint16(data[0:2])

	// 目标端口
	h.DestinationPort = binary.BigEndian.Uint16(data[2:4])

	// 序列号
	h.SequenceNumber = binary.BigEndian.Uint32(data[4:8])

	// 确认号
	h.Acknowledgment = binary.BigEndian.Uint32(data[8:12])

	// 数据偏移和标志
	h.DataOffset = data[12] >> 4
	h.Flags = data[12] & 0x3F

	// 窗口大小
	h.WindowSize = binary.BigEndian.Uint16(data[14:16])

	// 校验和
	h.Checksum = binary.BigEndian.Uint16(data[16:18])

	// 紧急指针
	h.UrgentPointer = binary.BigEndian.Uint16(data[18:20])

	// 选项
	optionsLength := int(h.DataOffset)*4 - TCPHeaderLength
	if optionsLength > 0 && len(data) >= TCPHeaderLength+optionsLength {
		h.Options = make([]byte, optionsLength)
		copy(h.Options, data[20:20+optionsLength])
	}

	return nil
}

// String 返回TCP头部的字符串表示
func (h *Header) String() string {
	flags := ""
	if h.Flags&FlagFIN != 0 {
		flags += "FIN "
	}
	if h.Flags&FlagSYN != 0 {
		flags += "SYN "
	}
	if h.Flags&FlagRST != 0 {
		flags += "RST "
	}
	if h.Flags&FlagPSH != 0 {
		flags += "PSH "
	}
	if h.Flags&FlagACK != 0 {
		flags += "ACK "
	}
	if h.Flags&FlagURG != 0 {
		flags += "URG "
	}

	return fmt.Sprintf("TCP Header: %d -> %d, Seq: %d, Ack: %d, Flags: [%s], Window: %d",
		h.SourcePort, h.DestinationPort, h.SequenceNumber, h.Acknowledgment, flags, h.WindowSize)
}

// HasFlag 检查是否包含指定标志
func (h *Header) HasFlag(flag uint8) bool {
	return (h.Flags & flag) != 0
}

// NewHeader 创建新的TCP头部
func NewHeader(srcPort, dstPort uint16, seqNum, ackNum uint32, flags uint8, windowSize uint16) *Header {
	return &Header{
		SourcePort:      srcPort,
		DestinationPort: dstPort,
		SequenceNumber:  seqNum,
		Acknowledgment:  ackNum,
		DataOffset:      5, // 20字节 = 5个32位字
		Flags:           flags,
		WindowSize:      windowSize,
		Checksum:        0, // 将由Marshal计算
		UrgentPointer:   0,
		Options:         nil,
	}
}

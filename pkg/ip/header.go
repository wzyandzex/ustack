package ip

import (
	"encoding/binary"
	"fmt"
	"net"
	"ustack/internal/utils"
)

const (
	// IP头部长度
	IPHeaderLength = 20

	// IP协议号
	ProtocolICMP = 1
	ProtocolTCP  = 6
	ProtocolUDP  = 17

	// IP标志
	FlagDF = 0x4000 // Don't Fragment
	FlagMF = 0x2000 // More Fragments
)

// Header IP头部结构
type Header struct {
	Version        uint8   // 版本号 (4)
	IHL            uint8   // 头部长度 (4)
	TOS            uint8   // 服务类型
	TotalLength    uint16  // 总长度
	Identification uint16  // 标识
	Flags          uint16  // 标志
	FragmentOffset uint16  // 片偏移
	TTL            uint8   // 生存时间
	Protocol       uint8   // 协议
	Checksum       uint16  // 校验和
	SourceIP       [4]byte // 源IP地址
	DestinationIP  [4]byte // 目标IP地址
}

// Marshal 将IP头部序列化为字节数组
func (h *Header) Marshal() ([]byte, error) {
	data := make([]byte, IPHeaderLength)

	// 版本和头部长度
	data[0] = (h.Version << 4) | h.IHL

	// 服务类型
	data[1] = h.TOS

	// 总长度
	binary.BigEndian.PutUint16(data[2:4], h.TotalLength)

	// 标识
	binary.BigEndian.PutUint16(data[4:6], h.Identification)

	// 标志和片偏移
	flagsAndOffset := (h.Flags & 0xE000) | (h.FragmentOffset & 0x1FFF)
	binary.BigEndian.PutUint16(data[6:8], flagsAndOffset)

	// TTL
	data[8] = h.TTL

	// 协议
	data[9] = h.Protocol

	// 校验和（先设为0）
	binary.BigEndian.PutUint16(data[10:12], 0)

	// 源IP地址
	copy(data[12:16], h.SourceIP[:])

	// 目标IP地址
	copy(data[16:20], h.DestinationIP[:])

	// 计算校验和
	h.Checksum = utils.CalculateChecksum(data)
	binary.BigEndian.PutUint16(data[10:12], h.Checksum)

	return data, nil
}

// Unmarshal 从字节数组解析IP头部
func (h *Header) Unmarshal(data []byte) error {
	if len(data) < IPHeaderLength {
		return fmt.Errorf("IP header too short: %d bytes", len(data))
	}

	// 版本和头部长度
	h.Version = data[0] >> 4
	h.IHL = data[0] & 0x0F

	// 服务类型
	h.TOS = data[1]

	// 总长度
	h.TotalLength = binary.BigEndian.Uint16(data[2:4])

	// 标识
	h.Identification = binary.BigEndian.Uint16(data[4:6])

	// 标志和片偏移
	flagsAndOffset := binary.BigEndian.Uint16(data[6:8])
	h.Flags = flagsAndOffset & 0xE000
	h.FragmentOffset = flagsAndOffset & 0x1FFF

	// TTL
	h.TTL = data[8]

	// 协议
	h.Protocol = data[9]

	// 校验和
	h.Checksum = binary.BigEndian.Uint16(data[10:12])

	// 源IP地址
	copy(h.SourceIP[:], data[12:16])

	// 目标IP地址
	copy(h.DestinationIP[:], data[16:20])

	return nil
}

// String 返回IP头部的字符串表示
func (h *Header) String() string {
	return fmt.Sprintf("IP Header: %s -> %s, Protocol: %d, TTL: %d, Length: %d",
		net.IP(h.SourceIP[:]).String(),
		net.IP(h.DestinationIP[:]).String(),
		h.Protocol,
		h.TTL,
		h.TotalLength)
}

// IsFragment 检查是否为分片
func (h *Header) IsFragment() bool {
	return (h.Flags&FlagMF) != 0 || h.FragmentOffset != 0
}

// IsFirstFragment 检查是否为第一个分片
func (h *Header) IsFirstFragment() bool {
	return h.FragmentOffset == 0
}

// NewHeader 创建新的IP头部
func NewHeader(srcIP, dstIP [4]byte, protocol uint8, totalLength uint16) *Header {
	return &Header{
		Version:        4,
		IHL:            5, // 20字节 = 5个32位字
		TOS:            0,
		TotalLength:    totalLength,
		Identification: 0, // 将由发送方设置
		Flags:          0,
		FragmentOffset: 0,
		TTL:            64,
		Protocol:       protocol,
		Checksum:       0, // 将由Marshal计算
		SourceIP:       srcIP,
		DestinationIP:  dstIP,
	}
}

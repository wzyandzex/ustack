package utils

import (
	"encoding/binary"
)

// CalculateChecksum 计算IP校验和
func CalculateChecksum(data []byte) uint16 {
	var sum uint32
	length := len(data)

	// 处理16位对齐的数据
	for i := 0; i < length-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}

	// 处理最后一个字节（如果长度为奇数）
	if length%2 == 1 {
		sum += uint32(data[length-1]) << 8
	}

	// 处理进位
	for sum > 0xffff {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	return uint16(^sum)
}

// CalculateTCPChecksum 计算TCP校验和
func CalculateTCPChecksum(tcpHeader, payload []byte, srcIP, dstIP []byte) uint16 {
	// 伪头部
	pseudoHeader := make([]byte, 12)
	copy(pseudoHeader[0:4], srcIP)
	copy(pseudoHeader[4:8], dstIP)
	pseudoHeader[8] = 0 // 保留字段
	pseudoHeader[9] = 6 // TCP协议号
	binary.BigEndian.PutUint16(pseudoHeader[10:12], uint16(len(tcpHeader)+len(payload)))

	// 组合数据
	data := make([]byte, 0, len(pseudoHeader)+len(tcpHeader)+len(payload))
	data = append(data, pseudoHeader...)
	data = append(data, tcpHeader...)
	data = append(data, payload...)

	return CalculateChecksum(data)
}

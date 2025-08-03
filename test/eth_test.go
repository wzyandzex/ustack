package test

import (
	"bytes"
	"testing"
	"ustack/pkg/eth"
)

func TestEthernetFrameMarshal(t *testing.T) {
	// 创建测试数据
	srcMAC := [6]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	dstMAC := [6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	payload := []byte("Hello, ustack!")

	frame := eth.NewFrame(srcMAC, dstMAC, eth.EtherTypeIPv4, payload)

	// 序列化
	data, err := frame.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal frame: %v", err)
	}

	// 验证长度
	expectedLength := eth.EthernetHeaderLength + len(payload)
	if len(data) != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, len(data))
	}

	// 验证目标MAC地址
	if !bytes.Equal(data[0:6], dstMAC[:]) {
		t.Errorf("Destination MAC mismatch")
	}

	// 验证源MAC地址
	if !bytes.Equal(data[6:12], srcMAC[:]) {
		t.Errorf("Source MAC mismatch")
	}

	// 验证以太网类型
	etherType := uint16(data[12])<<8 | uint16(data[13])
	if etherType != eth.EtherTypeIPv4 {
		t.Errorf("Expected ether type 0x%04x, got 0x%04x", eth.EtherTypeIPv4, etherType)
	}

	// 验证载荷
	if !bytes.Equal(data[14:], payload) {
		t.Errorf("Payload mismatch")
	}
}

func TestEthernetFrameUnmarshal(t *testing.T) {
	// 创建测试数据
	srcMAC := [6]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	dstMAC := [6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	payload := []byte("Hello, ustack!")

	// 手动构造帧数据
	data := make([]byte, eth.EthernetHeaderLength+len(payload))
	copy(data[0:6], dstMAC[:])
	copy(data[6:12], srcMAC[:])
	// 使用网络字节序设置以太网类型
	data[12] = 0x08 // 高字节
	data[13] = 0x00 // 低字节
	copy(data[14:], payload)

	// 解析
	frame := &eth.Frame{}
	err := frame.Unmarshal(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal frame: %v", err)
	}

	// 验证字段
	if frame.DestinationMAC != dstMAC {
		t.Errorf("Destination MAC mismatch")
	}

	if frame.SourceMAC != srcMAC {
		t.Errorf("Source MAC mismatch")
	}

	if frame.EtherType != eth.EtherTypeIPv4 {
		t.Errorf("Ether type mismatch")
	}

	if !bytes.Equal(frame.Payload, payload) {
		t.Errorf("Payload mismatch")
	}
}

func TestEthernetFrameRoundTrip(t *testing.T) {
	// 创建原始帧
	srcMAC := [6]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	dstMAC := [6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	payload := []byte("Test payload for round trip")

	originalFrame := eth.NewFrame(srcMAC, dstMAC, eth.EtherTypeIPv4, payload)

	// 序列化
	data, err := originalFrame.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// 反序列化
	newFrame := &eth.Frame{}
	err = newFrame.Unmarshal(data)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// 比较
	if originalFrame.DestinationMAC != newFrame.DestinationMAC {
		t.Errorf("Destination MAC mismatch after round trip")
	}

	if originalFrame.SourceMAC != newFrame.SourceMAC {
		t.Errorf("Source MAC mismatch after round trip")
	}

	if originalFrame.EtherType != newFrame.EtherType {
		t.Errorf("Ether type mismatch after round trip")
	}

	if !bytes.Equal(originalFrame.Payload, newFrame.Payload) {
		t.Errorf("Payload mismatch after round trip")
	}
}

func TestEthernetFrameBroadcast(t *testing.T) {
	broadcastMAC := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	srcMAC := [6]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}

	frame := eth.NewFrame(srcMAC, broadcastMAC, eth.EtherTypeIPv4, []byte("test"))

	if !frame.IsBroadcast() {
		t.Errorf("Frame should be broadcast")
	}

	if frame.IsMulticast() {
		t.Errorf("Broadcast frame should not be multicast")
	}
}

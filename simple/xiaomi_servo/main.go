package main

import (
	"context"
	"log"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

func main() {
	// 打开CAN网络接口
	conn, err := socketcan.DialContext(context.Background(), "udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("failed to dial CAN network interface: %v", err)
	}
	defer conn.Close()

	// 创建一个CAN帧，ID为127
	frame := can.Frame{
		ID:     127,
		Length: 8,
		Data:   [8]byte{1, 2, 3, 4, 5, 6, 7, 8}, // 示例数据
	}

	// 发送CAN帧
	tx := socketcan.NewTransmitter(conn)
	err = tx.TransmitFrame(context.Background(), frame)
	if err != nil {
		log.Fatalf("failed to dial CAN network interface: %v", err)
	}

	// 等待一段时间，以便可以查看发送效果
	time.Sleep(1 * time.Second)
}

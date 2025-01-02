package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go.einride.tech/can"
	"math"
)

type CyberGearController struct {
	Bus      interface{} // todo
	MotorId  int
	MainCNId int
	PMin     float64
	PMax     float64
	VMin     float64
	VMax     float64
	TMin     float64
	TMax     float64
	KpMin    float64 // 0.0 ~
	KpMax    float64 // ~ 500.0
	KdMin    float64 // 0.0 ~
	KdMax    float64 //    ~ 5.0
}

func NewController() CyberGearController {
	return CyberGearController{}
}

func (c *CyberGearController) EnableControlMode(torque, targetAngle, targetVelocity, Kp, Kd uint32) {
	// 生成29位的仲裁ID的组成部分
	// 也不知道干嘛的，先抄着
	var targetMin uint32 = 0.0
	var targetMax uint32 = 65535.0
	torqueMapped := c.linearMapping(torque, -12.0, 12.0, targetMin, targetMax)
	data2 := torqueMapped

	targetAngleMapped := c.linearMapping(targetAngle, -4*math.Pi, 4*math.Pi, targetMin, targetMax)
	targetVelocityMapped := c.linearMapping(targetVelocity, -30.0, 30.0, targetMin, targetMax)
	KpMapped := c.linearMapping(Kp, 0.0, 500.0, targetMin, targetMax)
	KdMapped := c.linearMapping(Kd, 0.0, 5.0, targetMin, targetMax)

	// 创建8字节的数据区
	// Create a bytes buffer and ensure byte order
	buf := new(bytes.Buffer)
	// The order of binary.Write calls matters and should match "HHHH" order in struct.pack
	err := binary.Write(buf, binary.LittleEndian, targetAngleMapped)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	err = binary.Write(buf, binary.LittleEndian, targetVelocityMapped)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	err = binary.Write(buf, binary.LittleEndian, KpMapped)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	err = binary.Write(buf, binary.LittleEndian, KdMapped)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	// Convert buffer to a slice of bytes
	data1 := buf.Bytes()

	// 发送CAN消息
	c.sendMessage(RunModeControlMode, data1, int(data2))
}

// sendMessage 发送CAN消息并接收响应。
//
// 参数:
// cmd_mode: 命令模式。
// data2: 数据区2。
// data1: 要发送的数据字节。
// timeout: 发送消息的超时时间(默认为200ms)。
//
// 返回:
// 一个元组, 包含接收到的消息数据和接收到的消息仲裁ID(如果有)。
func (c *CyberGearController) sendMessage(runModel RunMode, data1 []byte, data2 int) {
	arbitrationId := (int(runModel) << 24) | (data2 << 8) | c.MotorId
	message := can.Message(
		arbitrationId, data = data1, is_extended_id = True
	)

}

// linearMapping 不知道怎么映射的，先这样写吧
// 对输入值进行线性映射。
//
// 参数:
// value: 输入值。
// value_min: 输入值的最小界限。
// value_max: 输入值的最大界限。
// target_min: 输出值的最小界限。
// target_max: 输出值的最大界限。
//
// 返回:
// 映射后的值。
func (c *CyberGearController) linearMapping(value, valueMin, valueMax, targetMin, targetMax uint32) uint32 {
	return (value-valueMin)/(valueMax-valueMin)*(targetMax-targetMin) + targetMin
}

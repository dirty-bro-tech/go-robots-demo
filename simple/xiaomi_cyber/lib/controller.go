package lib

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"strings"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
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

	canConn net.Conn
}

func NewController() CyberGearController {
	return CyberGearController{}
}

func (c *CyberGearController) Init(ctx context.Context, network, address string) (err error) {
	c.canConn, err = socketcan.DialContext(context.Background(), network, address)
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		c.ReceiveMessage(ctx)
	}(ctx)

	return nil
}

// SetMode 设置控制模式。
func (c *CyberGearController) SetMode(mode RunMode) {
	switch expr {

	}

}

// SetRunMode 设置运控模式。
func (c *CyberGearController) SetRunMode() {
	c.SetMode(RunModeControlMode)
}

func (c *CyberGearController) WriteSingleParameter(paramName string, value int) {
	if f, ok := featureParams[paramName]; ok {
		/*
		 写入单个参数。
		  参数:
		  index: 参数索引。
		  value: 要设置的值。
		  format: 数据格式。
		  返回:
		  解析后的接收消息。
		*/

		encodeData, _ := c.formatData(f.value, f.Format, "encode")

	}
}

func (c *CyberGearController) SendCMDInControlMode(torque, targetAngle, targetVelocity, Kp, Kd uint32) {
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
	err = c.sendMessage(RunModeControlMode, data1, int(data2))
	if err != nil {
		panic(err)
	}
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
func (c *CyberGearController) sendMessage(runModel RunMode, data1 []byte, data2 int) (err error) {
	arbitrationId := (int(runModel) << 24) | (data2 << 8) | c.MotorId
	frame := can.Frame{
		ID:         uint32(arbitrationId),
		Data:       can.Data{data1[0], data1[1], data1[2], data1[3], data1[4], data1[5], data1[6], data1[7]},
		Length:     uint8(len(data1)),
		IsRemote:   false,
		IsExtended: true,
	}

	tx := socketcan.NewTransmitter(c.canConn)
	err = tx.TransmitFrame(context.Background(), frame)
	if err != nil {
		return fmt.Errorf("transmitter transmit frame failed: %v", err)
	}
	defer func(tx *socketcan.Transmitter) {
		errIn := tx.Close()
		if errIn != nil {
			fmt.Println("transmitter close failed:", err)
		}
	}(tx)

	return nil
}

// ReceiveMessage 接收CAN消息。
//
// 参数:
// timeout: 接收消息的超时时间(默认为200ms)。
//
// 返回:
// 一个元组, 包含接收到的消息数据和接收到的消息仲裁ID(如果有)。
func (c *CyberGearController) ReceiveMessage(ctx context.Context) {
	// 接收CAN消息
	recv := socketcan.NewReceiver(c.canConn)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, exiting ReceiveMessage loop.")
			return
		default:
			if recv.Receive() {
				frame := recv.Frame()
				fmt.Println(frame.String())
			}
		}
	}
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

// formatData handles encoding or decoding of data based on the provided format.
func (c *CyberGearController) formatData(data []byte, format string, opType string) ([]interface{}, error) {
	formatList := strings.Fields(format)
	var result []interface{}

	if opType == "decode" {
		reader := bytes.NewReader(data)
		for _, f := range formatList {
			switch f {
			case "f":
				var value float32
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "u16":
				var value uint16
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "s16":
				var value int16
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "u32":
				var value uint32
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "s32":
				var value int32
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "u8":
				var value uint8
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			case "s8":
				var value int8
				if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				result = append(result, value)
			default:
				log.Printf("unknown format in FormatData(): %s", f)
				return nil, fmt.Errorf("unknown format: %s", f)
			}
		}
		return result, nil
	} else if opType == "encode" {
		var buf bytes.Buffer
		for i, f := range formatList {
			switch f {
			case "f":
				if err := binary.Write(&buf, binary.LittleEndian, float32(data[i])); err != nil {
					return nil, err
				}
			case "u16", "s16", "u32", "s32", "u8", "s8":
				if err := binary.Write(&buf, binary.LittleEndian, data[i]); err != nil {
					return nil, err
				}
			default:
				log.Printf("unknown format in FormatData(): %s", f)
				return nil, fmt.Errorf("unknown format: %s", f)
			}
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("invalid operation type: %s", opType)
}

package lib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"

	"encoding/binary"
	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

type CyberGearController struct {
	Bus       interface{} // todo
	MotorId   int
	MainCanId int
	PMin      float64
	PMax      float64
	VMin      float64
	VMax      float64
	TMin      float64
	TMax      float64
	KpMin     float64 // 0.0 ~
	KpMax     float64 // ~ 500.0
	KdMin     float64 // 0.0 ~
	KdMax     float64 //    ~ 5.0

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
	switch mode {
	case RunModeControlMode:
		c.WriteSingleParameter(featureParams["run_mode"].Index, "run_mode", mode.Mode().Value())
	case RunModePositionMode:
		log.Printf("WARNING: SetMode: unknown run mode: %v", mode)
	case RunModeSpeedMode:
		log.Printf("WARNING: SetMode: unknown run mode: %v", mode)
	case RunModeCurrentMode:
	default:
		log.Printf("WARNING: SetMode: unknown run mode: %v", mode)
		return
	}
}

// SetRunMode 设置运控模式。
func (c *CyberGearController) SetRunMode() {
	c.SetMode(RunModeControlMode)
}

// Disable stops motor
func (c *CyberGearController) Disable() {
	c.clearCanRxBuffer()

	err := c.sendMessage(CMDModeMotorStop.Mode(), c.MainCanId, []byte{0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		log.Println(err)
	}
}

// Enable enables motor
//
//	使能运行电机。
func (c *CyberGearController) Enable() {
	c.clearCanRxBuffer()

	err := c.sendMessage(CMDModeMotorEnable.Mode(), c.MainCanId, []byte{})
	if err != nil {
		log.Println(err)
	}
}

func (c *CyberGearController) WriteSingleParameter(index int, paramName string, value int) {
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

		encodedData, _ := c.formatData([]byte{byte(value)}, f.Format, "encode")
		// Create a bytes buffer and write the index as little-endian bytes
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, index)
		if err != nil {
			panic(err)
		}

		// Convert buffer bytes to a slice
		indexBytes := buf.Bytes()

		// Concatenate indexBytes and encodedData
		data1 := append(indexBytes, encodedData...)
		c.clearCanRxBuffer()

		err = c.sendMessage(CMDModeSingleParamWrite.Mode(), c.MainCanId, data1)
		if err != nil {
			log.Println(err)
		}
	}
}

// SendCMDInControlMode 运控模式下发送电机控制指令。
//
//	参数:
//	torque: 扭矩。
//	target_angle: 目标角度。
//	target_velocity: 目标速度。
//	Kp: 比例增益。
//	Kd: 导数增益。
//
//	返回:
//	解析后的接收消息。
func (c *CyberGearController) SendCMDInControlMode(torque, targetAngle, targetVelocity, Kp, Kd uint32) {
	// 生成29位的仲裁ID的组成部分
	// 也不知道干嘛的，先抄着
	var targetMin float32 = 0.0
	var targetMax float32 = 65535.0
	torqueMapped := c.linearMapping(float32(torque), -12, 12, targetMin, targetMax)
	data2 := torqueMapped

	targetAngleMapped := c.linearMapping(float32(targetAngle), -4*math.Pi, 4*math.Pi, targetMin, targetMax)
	targetVelocityMapped := c.linearMapping(float32(targetVelocity), -30.0, 30.0, targetMin, targetMax)
	KpMapped := c.linearMapping(float32(Kp), 0.0, 500.0, targetMin, targetMax)
	KdMapped := c.linearMapping(float32(Kd), 0.0, 5.0, targetMin, targetMax)

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
	err = c.sendMessage(RunModeControlMode.Mode(), int(data2), data1)
	if err != nil {
		panic(err)
	}
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
func (c *CyberGearController) linearMapping(value, valueMin, valueMax, targetMin, targetMax float32) uint32 {
	return uint32((value-valueMin)/(valueMax-valueMin)*(targetMax-targetMin) + targetMin)
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
func (c *CyberGearController) sendMessage(mode Mode, data2 int, data1 []byte) (err error) {
	arbitrationId := (int(mode) << 24) | (data2 << 8) | c.MotorId
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

// 解析接收到的CAN消息。
//
// 参数:
// data: 接收到的数据。
// arbitration_id: 接收到的消息的仲裁ID。
//
// 返回:
// 一个元组, 包含电机的CAN ID、位置(rad)、速度(rad/s)、力矩(Nm)、温度(摄氏度)。
func (c *CyberGearController) parseReceivedMsg(data []byte, arbitrationID uint32) (uint8, float64, float64, float64, float64) {
	if data != nil {
		log.Printf("Received message with ID %x", arbitrationID)

		// Parse Motor CAN ID
		motorCanID := uint8((arbitrationID >> 8) & 0xFF)

		pos := uintToFloat(binary.BigEndian.Uint16(data[0:2]), c.PMin, c.PMax, twoBytesBits)
		vel := uintToFloat(binary.BigEndian.Uint16(data[2:4]), c.VMin, c.VMax, twoBytesBits)
		torque := uintToFloat(binary.BigEndian.Uint16(data[4:6]), c.TMin, c.TMax, twoBytesBits)

		// Parse temperature data
		temperatureRaw := binary.BigEndian.Uint16(data[6:8])
		temperatureCelsius := float64(temperatureRaw) / 10.0

		log.Printf("Motor CAN ID: %d, pos: %.2f rad, vel: %.2f rad/s, torque: %.2f Nm, temperature: %.1f °C",
			motorCanID, pos, vel, torque, temperatureCelsius)

		return motorCanID, pos, vel, torque, temperatureCelsius
	} else {
		log.Println("No message received within the timeout period.")
		return 0, 0, 0, 0, 0
	}
}

func uintToFloat(value uint16, min, max float64, bits uint) float64 {
	return min + (float64(value)/math.Pow(2, float64(bits)))*(max-min)
}

func (c *CyberGearController) clearCanRxBuffer() {
}

// formatData encodes or decodes data to/from a byte slice depending on the specified format and operation type.
func (c *CyberGearController) formatData(data []byte, format string, dataType string) ([]byte, error) {
	formatList := splitWhitespace(format)
	buf := bytes.NewBuffer(nil)

	if dataType == "decode" {
		inputBuffer := bytes.NewBuffer(data)
		for _, f := range formatList {
			switch f {
			case "f": // float32, 4 bytes
				var value float32
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "u16": // uint16, 2 bytes
				var value uint16
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "s16": // int16, 2 bytes
				var value int16
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "u32": // uint32, 4 bytes
				var value uint32
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "s32": // int32, 4 bytes
				var value int32
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "u8": // uint8, 1 byte
				var value uint8
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			case "s8": // int8, 1 byte
				var value int8
				if err := binary.Read(inputBuffer, binary.LittleEndian, &value); err != nil {
					return nil, err
				}
				binary.Write(buf, binary.LittleEndian, value)
			default:
				return nil, errors.New("unknown format: " + f)
			}
		}
		return buf.Bytes(), nil
	} else if dataType == "encode" {
		// Assuming `data` is a properly formatted byte slice for encoding based on `format`
		buf.Write(data)
		return buf.Bytes(), nil
	}

	return nil, errors.New("invalid type specified, must be 'encode' or 'decode'")
}

// splitWhitespace splits a string by whitespace and returns a slice of strings.
func splitWhitespace(s string) []string {
	fields := bytes.Fields([]byte(s))
	result := make([]string, len(fields))
	for i, field := range fields {
		result[i] = string(field)
	}
	return result
}

func uint32ToBytes(num uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, num)
	return buf.Bytes()
}

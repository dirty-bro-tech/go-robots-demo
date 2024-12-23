package main

import (
	"machine"
	"time"
)

const (
	startPin           = machine.D9
	breakPin           = machine.D9
	changeDirectionPin = machine.D3
	speedFeedbackPin   = machine.D6
)

type MotorPMWI interface {
	Configure(config machine.PWMConfig) error
	Channel(pin machine.Pin) (channel uint8, err error)
	Top() uint32
	Set(channel uint8, value uint32)
}

type MotorPWM struct {
}

// MotorServo 无刷电机
type MotorServo struct {
	// 启动Channel
	startChannel MotorServoChannel
	// 切换方向Channel
	changeDirectionChannel MotorServoChannel
	// 刹车Channel
	breakChannel MotorServoChannel
	// 速度反馈Channel
	speedFeedbackChannel MotorServoChannel
}

func (m *MotorServo) Init() error {
	if err := m.initStartChannel(); err != nil {
		return err
	}
	if err := m.initBreakChannel(); err != nil {
		return err
	}

	return nil
}

const pwmPeriod = 20e6 // 20ms
func (m *MotorServo) initStartChannel() error {
	// 定义 PWM 周期
	period := pwmPeriod
	// 获取 PWM 设备
	pwm := machine.Timer1
	// 配置 PWM 设备
	err := pwm.Configure(machine.PWMConfig{Period: uint64(period)})
	if err != nil {
		println("start PWM err:", err)
		return err
	}

	m.startChannel.channel, err = pwm.Channel(m.startChannel.pin)
	if err != nil {
		println("start chan err:", err)
		return err
	}

	m.startChannel.pwm = pwm

	return nil
}

func (m *MotorServo) initBreakChannel() error {
	// 定义 PWM 周期
	period := pwmPeriod
	// 获取 PWM 设备
	pwm := machine.Timer1
	// 配置 PWM 设备
	err := pwm.Configure(machine.PWMConfig{Period: uint64(period)})
	if err != nil {
		println("brk PWM err:", err)
		return err
	}

	m.breakChannel.channel, err = pwm.Channel(m.breakChannel.pin)
	if err != nil {
		println("brk chan err:", err)
		return err
	}

	m.breakChannel.pwm = pwm

	return nil
}

func (m *MotorServo) start2() {
	println("start")
	// microseconds := 80*1000/180 + 1000
	// value := uint64(m.startChannel.pwm.Top()) * uint64(microseconds) / (pwmPeriod / 1000)

	m.startChannel.pwm.Set(m.startChannel.channel, m.startChannel.pwm.Top()/32)
	println("start end")
}

func (m *MotorServo) start() {
	println("start")
	startPin2 := machine.D9 // BRK 引脚
	startPin2.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// 释放 BRK 信号（低电平）
	startPin2.High()
	println("start end")
}

func (m *MotorServo) brk() {
	println("brk")

	// 配置 BRK 信号引脚
	brkPin := machine.D8 // BRK 引脚
	brkPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// 释放 BRK 信号（低电平）
	brkPin.Low()
	// 模拟电机停止
	println("brk end")
}

func (m *MotorServo) changeDirection() {
	println("change")

	// 配置 BRK 信号引脚
	changePin := machine.D3 // BRK 引脚
	changePin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// 释放 BRK 信号（低电平）
	changePin.Low()
	// 模拟电机停止
	println("change end")
}

type MotorServoChannel struct {
	pin     machine.Pin
	pwm     machine.PWM
	channel uint8
}

func NewMotorServo() *MotorServo {
	servo := &MotorServo{}
	// 初始化
	servo.startChannel = MotorServoChannel{
		pin: startPin,
	}
	servo.changeDirectionChannel = MotorServoChannel{
		pin: changeDirectionPin,
	}
	servo.breakChannel = MotorServoChannel{
		pin: breakPin,
	}
	servo.speedFeedbackChannel = MotorServoChannel{
		pin: speedFeedbackPin,
	}

	return servo
}

func main() {
	motor := NewMotorServo()
	if err := motor.Init(); err != nil {
		return
	}

	// 初始化 UART
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
	})

	for {
		motor.start()
		time.Sleep(60 * time.Millisecond)
		motor.changeDirection()
		time.Sleep(60 * time.Millisecond)
		motor.brk()
	}
	// motor.brk()
}

func strToInt(s string) (int, bool) {
	result := 0
	sign := 1
	start := 0

	if len(s) == 0 {
		return 0, false
	}

	// 处理负号
	if s[0] == '-' {
		sign = -1
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
		result = result*10 + int(s[i]-'0')
	}

	return sign * result, true
}

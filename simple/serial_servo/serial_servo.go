package main

import (
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/servo"
)

// Configuration for the Arduino Uno.
// Please change the PWM and pin if you want to try this example on a different
// board.
var (
	pwm = machine.Timer1
	pin = machine.D9
)

func main() {
	// 初始化 UART
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
	})

	s, err := servo.New(pwm, pin)
	if err != nil {
		for {
			println("could not configure servo")
			time.Sleep(time.Second)
		}
	}

	for {
		// 读取 UART 输入
		if uart.Buffered() > 0 {
			input := make([]byte, 10)
			n, _ := uart.Read(input)
			cmd := string(input[:n])

			// 转换输入为角度
			angle, _ := strToInt(cmd)
			if angle >= 0 && angle <= 180 {
				// 计算脉冲宽度
				// pulseWidth := angleToPulseWidth(angle)
				s.SetAngle(angle)
				// 打印调试信息
				uart.Write([]byte("Angle set to: " + strconv.Itoa(angle) + "\n"))
			} else {
				uart.Write([]byte("Invalid angle. Please enter a value between 0 and 180.\n"))
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
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

// 将角度转换为脉冲宽度
func angleToPulseWidth(angle int) int {
	return 500 + (angle * 2000 / 180)
}

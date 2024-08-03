package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyUSB0")
	servo := gpio.NewServoDriver(firmataAdaptor, "9")

	work := func() {
		gobot.Every(1*time.Second, func() {
			_ = servo.Move(10)
		})
	}

	robot := gobot.NewRobot("toggle-every-one-second",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	err := robot.Start()
	if err != nil {
		panic(err)
	}
}

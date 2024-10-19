package main

import "github.com/brutella/can"

func main() {
	bus, err := can.NewBusForInterfaceWithName("can0")
	if err != nil {
		panic(err)
	}
	bus.ConnectAndPublish()

	frm := can.Frame{
		ID:     0x701,
		Length: 1,
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   [8]uint8{0x05},
	}

	bus.Publish(frm)
}

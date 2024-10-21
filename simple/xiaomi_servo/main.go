package main

import "github.com/brutella/can"

func main() {
	bus, err := can.NewBusForInterfaceWithName("can1")
	if err != nil {
		panic(err)
	}

	err = bus.ConnectAndPublish()
	if err != nil {
		panic(err)
	}

	frm := can.Frame{
		ID:     8,
		Length: 1,
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   [8]uint8{0x05},
	}

	err = bus.Publish(frm)
	if err != nil {
		panic(err)
	}
}

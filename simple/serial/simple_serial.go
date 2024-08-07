package main

import (
	"machine"
	"time"
)

func main() {
	println("Reading from the serial port...")

	for {
		c, err := machine.Serial.ReadByte()
		if err == nil {
			if c < 32 {
				// Convert nonprintable control characters to
				// ^A, ^B, etc.
				machine.Serial.WriteByte('^')
				machine.Serial.WriteByte(c + '@')
			} else if c >= 127 {
				// Anything equal or above ASCII 127, print ^?.
				machine.Serial.WriteByte('^')
				machine.Serial.WriteByte('?')
			} else {
				// Echo the printable character back to the
				// host computer.
				machine.Serial.WriteByte(c)
			}
		}

		// This assumes that the input is coming from a keyboard
		// so checking 120 times per second is sufficient. But if
		// the data comes from another processor, the port can
		// theoretically receive as much as 11000 bytes/second
		// (115200 baud). This delay can be removed and the
		// Serial.Read() method can be used to retrieve
		// multiple bytes from the receive buffer for each
		// iteration.
		time.Sleep(time.Millisecond * 8)
	}
}

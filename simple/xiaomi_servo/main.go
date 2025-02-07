package main

import (
	"context"
	"github.com/chip-ai-labs/go-robots-demo/simple/xiaomi_cyber/lib"
)

func main() {
	c := lib.NewController()
	c.Init(context.Background(), "can", "can1")
	c.SetRunMode()
}

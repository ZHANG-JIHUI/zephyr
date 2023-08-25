package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/remote"
)

func main() {

	engine := actor.NewEngine()
	remoter := remote.NewRemoter(engine, ":8880")
	engine.WithRemoter(remoter)

	engine.SpawnFunc(func(ctx *actor.Context) {
		switch message := ctx.Message().(type) {
		case actor.Started:
		case actor.Stopped:
		case *actor.PID:
			ctx.Send(message, &network.NetPacketMessage{
				Id:   1000,
				Data: []byte("hello world"),
			})
		default:
		}
	}, "master")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

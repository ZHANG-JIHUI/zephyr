package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ZHANG-JIHUI/zephyr/cluster/gateway"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

func main() {
	logger := log.NewLog(log.WithLogDir("./logs", "./logs"), log.WithRunMode(log.RunModeDev, nil))
	log.SetLogger(logger)

	engine := actor.NewEngine()

	pid := engine.Spawn(gateway.NewGateway().Producer(), "gateway")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	var wg sync.WaitGroup
	engine.Poison(pid, &wg)
	wg.Wait()
}

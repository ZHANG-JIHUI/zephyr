package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ZHANG-JIHUI/zephyr/network"
	"github.com/ZHANG-JIHUI/zephyr/network/tcp"
	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

func main() {
	logger := log.NewLog(log.WithLogDir("./logs", "./logs"), log.WithRunMode(log.RunModeDev, nil))
	log.SetLogger(logger)

	server := tcp.NewServer(":9999")
	server.RegStartEvent(func(srv network.Server) {
		log.Info("tcp server started", log.String("addr", srv.Addr()), log.Any("pid", srv.PID()))
	})
	server.RegStopEvent(func(srv network.Server) {
		log.Info("tcp server stopped", log.String("addr", srv.Addr()), log.Any("pid", srv.PID()))
	})
	server.RegConnectEvent(func(srv network.Server, conn network.Conn) {
		log.Info("client connected", log.Any("connection pid", conn.PID()))
	})
	server.RegDisconnectEvent(func(srv network.Server, conn network.Conn) {
		log.Info("client disconnected", log.Any("connection pid", conn.PID()))
	})
	server.RegReceiveEvent(func(srv network.Server, conn network.Conn, msg []byte, typ int) {
		log.Info("received message", log.Any("pid", conn.PID()), log.ByteString("msg", msg))
	})

	engine := actor.NewEngine()
	pid := engine.Spawn(server.Producer(), "server")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	var wg sync.WaitGroup
	engine.Poison(pid, &wg)
	wg.Wait()
}

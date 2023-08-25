package main

import (
	"net"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}

	go func() {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		log.Info("received", log.ByteString("msg", buf[:n]))
	}()

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()
	for range ticker.C {
		if _, err := conn.Write([]byte("hello")); err != nil {
			return
		}
	}
}

package main

import (
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for range ticker.C {
		if _, err := conn.Write([]byte("hello")); err != nil {
			return
		}
	}
}

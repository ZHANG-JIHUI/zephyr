package network

import (
	"net"

	"github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"
)

type Conn interface {
	actor.Base
	ID() string
	RemoteAddr() net.Addr
	IP() string
	Close() error
}

package network

import "github.com/ZHANG-JIHUI/zephyr/third_party/hollywood/actor"

type Server interface {
	actor.Base
	Addr() string
	Protocol() string
	Start() error
	Stop() error
	RegStartEvent(handler StartEventHandler)
	RegStopEvent(handler StopEventHandler)
	RegConnectEvent(handler ConnectEventHandler)
	RegReceiveEvent(handler ReceiveEventHandler)
	RegDisconnectEvent(handler DisconnectEventHandler)
	OnStartEvent()
	OnStopEvent()
	OnConnectEvent(conn Conn)
	OnDisconnectEvent(conn Conn)
	OnReceiveEvent(conn Conn, msg []byte, typ int)
}

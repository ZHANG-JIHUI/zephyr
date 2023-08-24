package network

type (
	StartEventHandler      func(srv Server)
	StopEventHandler       func(srv Server)
	ConnectEventHandler    func(srv Server, conn Conn)
	DisconnectEventHandler func(srv Server, conn Conn)
	ReceiveEventHandler    func(srv Server, conn Conn, msg []byte, typ int)
)

type Event struct {
	Server
	StartEventHandler      StartEventHandler
	StopEventHandler       StopEventHandler
	ConnectEventHandler    ConnectEventHandler
	DisconnectEventHandler DisconnectEventHandler
	ReceiveEventHandler    ReceiveEventHandler
}

func (slf *Event) RegStartEvent(handler StartEventHandler) {
	slf.StartEventHandler = handler
}

func (slf *Event) OnStartEvent() {
	slf.StartEventHandler(slf.Server)
}

func (slf *Event) RegStopEvent(handler StopEventHandler) {
	slf.StopEventHandler = handler
}

func (slf *Event) OnStopEvent() {
	slf.StopEventHandler(slf.Server)
}

func (slf *Event) RegConnectEvent(handler ConnectEventHandler) {
	slf.ConnectEventHandler = handler
}

func (slf *Event) OnConnectEvent(conn Conn) {
	slf.ConnectEventHandler(slf.Server, conn)
}

func (slf *Event) RegDisconnectEvent(handler DisconnectEventHandler) {
	slf.DisconnectEventHandler = handler
}

func (slf *Event) OnDisconnectEvent(conn Conn) {
	slf.DisconnectEventHandler(slf.Server, conn)
}

func (slf *Event) RegReceiveEvent(handler ReceiveEventHandler) {
	slf.ReceiveEventHandler = handler
}

func (slf *Event) OnReceiveEvent(conn Conn, msg []byte, typ int) {
	slf.ReceiveEventHandler(slf.Server, conn, msg, typ)
}

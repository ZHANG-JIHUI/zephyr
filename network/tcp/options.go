package tcp

import "github.com/ZHANG-JIHUI/zephyr/network"

type options struct {
	runMode    network.RunMode
	addr       string
	protocol   string
	maxMsgLen  int
	maxConnNum int
}

type Option func(*options)

func defaultOptions(addr string) *options {
	return &options{
		runMode:  network.RunModeDev,
		addr:     addr,
		protocol: "tcp",
	}
}

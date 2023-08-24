package network

import "github.com/ZHANG-JIHUI/zephyr/tools/log"

type RunMode = log.RunMode

const (
	RunModeDev  RunMode = log.RunModeDev
	RunModeTest RunMode = log.RunModeTest
	RunModeProd RunMode = log.RunModeProd
)

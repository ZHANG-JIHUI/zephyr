package log

type RunMode = uint8

const (
	RunModeDev  RunMode = iota
	RunModeTest RunMode = iota
	RunModeProd RunMode = iota
)

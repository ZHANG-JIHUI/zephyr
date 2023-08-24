package log

import "testing"

func TestLogger(t *testing.T) {
	logger := NewLog(WithLogDir("./logs", "./logs"), WithRunMode(RunModeDev, nil))
	SetLogger(logger)

	logger.Info("info msg")
	Info("info msg")
}

package log

import (
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

var logger Logger = NewLog()

type Logger interface {
	ants.Logger
	logging.Logger
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	DPanicf(format string, args ...any)
	Panicf(format string, args ...any)
	Fatalf(format string, args ...any)
}

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func DPanicf(format string, args ...any) {
	logger.DPanicf(format, args...)
}

func Panicf(format string, args ...any) {
	logger.Panicf(format, args...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

func SetLogger(log Logger) {
	logger = log
}

func GetLogger() Logger {
	return logger
}

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level
type LevelEnablerFunc = zap.LevelEnablerFunc

const (
	DebugLevel  Level = zapcore.DebugLevel
	InfoLevel   Level = zapcore.InfoLevel
	WarnLevel   Level = zapcore.WarnLevel
	ErrorLevel  Level = zapcore.ErrorLevel
	DPanicLevel Level = zapcore.DPanicLevel
	PanicLevel  Level = zapcore.PanicLevel
	FatalLevel  Level = zapcore.FatalLevel
)

var (
	levels                = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, FatalLevel}
	defaultLevelPartition = map[Level]func() LevelEnablerFunc{
		DebugLevel:  DebugLevelPartition,
		InfoLevel:   InfoLevelPartition,
		WarnLevel:   WarnLevelPartition,
		ErrorLevel:  ErrorLevelPartition,
		DPanicLevel: DPanicLevelPartition,
		PanicLevel:  PanicLevelPartition,
		FatalLevel:  FatalLevelPartition,
	}
)

func Levels() []Level {
	return levels
}

func DebugLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == DebugLevel
	}
}

func InfoLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == InfoLevel
	}
}

func WarnLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == WarnLevel
	}
}

func ErrorLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == ErrorLevel
	}
}

func DPanicLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == DPanicLevel
	}
}

func PanicLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == PanicLevel
	}
}

func FatalLevelPartition() LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level == FatalLevel
	}
}

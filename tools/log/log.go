package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Log struct {
	zap            *zap.Logger
	sugar          *zap.SugaredLogger
	filename       func(level Level) string
	rotateFilename func(level Level) string
	rotateOptions  []rotatelogs.Option
	levelPartition map[Level]func() LevelEnablerFunc
	cores          []Core
	mode           RunMode
	logDir         string
	rotateLogDir   string
}

func NewLog(options ...Option) *Log {
	log := &Log{
		filename: func(level Level) string {
			return fmt.Sprintf("%s.log", level.String())
		},
		rotateFilename: func(level Level) string {
			return strings.Join([]string{level.String(), "%Y%m%d.log"}, ".")
		},
		levelPartition: defaultLevelPartition,
	}

	for _, option := range options {
		option(log)
	}

	if len(log.rotateOptions) == 0 {
		log.rotateOptions = []rotatelogs.Option{
			rotatelogs.WithMaxAge(time.Hour * 24 * 7),
			rotatelogs.WithRotationTime(time.Hour * 24),
		}
	}

	if len(log.cores) == 0 {
		var encoder = NewEncoder()

		switch log.mode {
		case RunModeDev:
			var partition LevelEnablerFunc = func(lvl Level) bool {
				return true
			}
			log.cores = append(log.cores, zapcore.NewCore(encoder, os.Stdout, partition))
		case RunModeTest:
			infoRotate, err := rotatelogs.New(
				filepath.Join(log.rotateLogDir, log.rotateFilename(InfoLevel)),
				append([]rotatelogs.Option{rotatelogs.WithLinkName(filepath.Join(log.logDir, log.filename(InfoLevel)))}, log.rotateOptions...)...,
			)
			if err != nil {
				panic(err)
			}
			errRotate, err := rotatelogs.New(
				filepath.Join(log.rotateLogDir, log.rotateFilename(ErrorLevel)),
				append([]rotatelogs.Option{rotatelogs.WithLinkName(filepath.Join(log.logDir, log.filename(ErrorLevel)))}, log.rotateOptions...)...,
			)
			if err != nil {
				panic(err)
			}
			if log.logDir != "" {
				log.cores = append(log.cores, zapcore.NewCore(encoder, zapcore.AddSync(infoRotate), LevelEnablerFunc(func(lvl Level) bool { return lvl < ErrorLevel })))
				log.cores = append(log.cores, zapcore.NewCore(encoder, zapcore.AddSync(errRotate), LevelEnablerFunc(func(lvl Level) bool { return lvl >= ErrorLevel })))
				log.cores = append(log.cores, zapcore.NewCore(encoder, os.Stdout, LevelEnablerFunc(func(lvl Level) bool { return lvl < ErrorLevel })))
				log.cores = append(log.cores, zapcore.NewCore(encoder, os.Stdout, LevelEnablerFunc(func(lvl Level) bool { return lvl >= ErrorLevel })))
			}
		case RunModeProd:
			infoRotate, err := rotatelogs.New(
				filepath.Join(log.rotateLogDir, log.rotateFilename(InfoLevel)),
				append([]rotatelogs.Option{rotatelogs.WithLinkName(filepath.Join(log.logDir, log.filename(InfoLevel)))}, log.rotateOptions...)...,
			)
			if err != nil {
				panic(err)
			}
			errRotate, err := rotatelogs.New(
				filepath.Join(log.rotateLogDir, log.rotateFilename(ErrorLevel)),
				append([]rotatelogs.Option{rotatelogs.WithLinkName(filepath.Join(log.logDir, log.filename(ErrorLevel)))}, log.rotateOptions...)...,
			)
			if err != nil {
				panic(err)
			}
			if log.logDir != "" {
				log.cores = append(log.cores, zapcore.NewCore(encoder, zapcore.AddSync(infoRotate), LevelEnablerFunc(func(lvl Level) bool { return lvl >= InfoLevel && lvl <= WarnLevel })))
				log.cores = append(log.cores, zapcore.NewCore(encoder, zapcore.AddSync(errRotate), LevelEnablerFunc(func(lvl Level) bool { return lvl >= ErrorLevel })))
			}
		}
	}

	log.zap = zap.New(zapcore.NewTee(log.cores...), zap.AddCaller(), zap.AddCallerSkip(2))
	log.sugar = log.zap.Sugar()
	return log
}

func (slf *Log) Debugf(format string, args ...any) {
	slf.sugar.Debugf(format, args...)
}

func (slf *Log) Infof(format string, args ...any) {
	slf.sugar.Infof(format, args...)
}

func (slf *Log) Warnf(format string, args ...any) {
	slf.sugar.Warnf(format, args...)
}

func (slf *Log) Errorf(format string, args ...any) {
	slf.sugar.Errorf(format, args...)
}

func (slf *Log) DPanicf(format string, args ...any) {
	slf.sugar.DPanicf(format, args...)
}

func (slf *Log) Panicf(format string, args ...any) {
	slf.sugar.Panicf(format, args...)
}

func (slf *Log) Fatalf(format string, args ...any) {
	slf.sugar.Fatalf(format, args...)
}

func (slf *Log) Printf(format string, args ...any) {
	slf.sugar.Infof(format, args...)
}

func (slf *Log) Debug(msg string, fields ...Field) {
	slf.zap.Debug(msg, fields...)
}

func (slf *Log) Info(msg string, fields ...Field) {
	slf.zap.Info(msg, fields...)
}

func (slf *Log) Warn(msg string, fields ...Field) {
	slf.zap.Warn(msg, fields...)
}

func (slf *Log) Error(msg string, fields ...Field) {
	slf.zap.Error(msg, fields...)
}

// DPanic 仅在开发模式下Panic
func (slf *Log) DPanic(msg string, fields ...Field) {
	slf.zap.DPanic(msg, fields...)
}

func (slf *Log) Panic(msg string, fields ...Field) {
	slf.zap.Panic(msg, fields...)
}

func (slf *Log) Fatal(msg string, fields ...Field) {
	slf.zap.Fatal(msg, fields...)
}

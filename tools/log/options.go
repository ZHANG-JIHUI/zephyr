package log

import rotatelogs "github.com/lestrrat-go/file-rotatelogs"

type Option func(log *Log)

// WithRunMode 设置运行模式，默认为：RunModeDev
func WithRunMode(mode RunMode, handle func() Core) Option {
	return func(log *Log) {
		log.mode = mode
		if handle != nil {
			log.cores = append(log.cores, handle())
		}
	}
}

// WithFilename 设置日志文件名，默认为：{level}.log
func WithFilename(filename func(level Level) string) Option {
	return func(log *Log) {
		log.filename = filename
	}
}

// WithRotateFilename 设置日志分割文件名，默认为：{level}.%Y%m%d.log
func WithRotateFilename(filename func(level Level) string) Option {
	return func(log *Log) {
		log.rotateFilename = filename
	}
}

// WithRotateOption 设置日志分割选项，默认选项为：WithMaxAge(time.Hour*24*7)，WithRotationTime(time.Hour*24)
func WithRotateOption(options ...rotatelogs.Option) Option {
	return func(log *Log) {
		log.rotateOptions = options
	}
}

// WithLogDir 设置日志存储目录，若为空则不会存储文件
func WithLogDir(logDir, rotateLogDir string) Option {
	return func(log *Log) {
		log.logDir = logDir
		log.rotateLogDir = rotateLogDir
	}
}

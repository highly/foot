package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Level int

const (
	_ Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var defaultLogger *logger

type logger struct {
	zl    *zap.Logger
	level Level
}

func New(logLevel int) {
	options := []zap.Option{
		zap.AddCallerSkip(1),
	}
	cfg := zap.NewProductionConfig()
	log, err := cfg.Build(options...)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Logger:fail to initialize logger, msg:", err)
		os.Exit(-1)
	}
	defaultLogger = &logger{zl: log, level: Level(logLevel)}
}

func Error(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= ErrorLevel {
		defaultLogger.zl.Error(msg, fields...)
	}
}

func Errorf(template string, args ...interface{}) {
	if defaultLogger.level >= ErrorLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zl.Error(msg)
	}
}

func Warn(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= WarnLevel {
		defaultLogger.zl.Warn(msg, fields...)
	}
}

func Warnf(template string, args ...interface{}) {
	if defaultLogger.level >= WarnLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zl.Warn(msg)
	}
}

func Info(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= InfoLevel {
		defaultLogger.zl.Info(msg, fields...)
	}
}

func Infof(template string, args ...interface{}) {
	if defaultLogger.level >= InfoLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zl.Info(msg)
	}
}

func Debug(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= DebugLevel {
		defaultLogger.zl.Debug(msg, fields...)
	}
}

func Debugf(template string, args ...interface{}) {
	if defaultLogger.level >= DebugLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zl.Debug(msg)
	}
}

func Sync() error {
	return defaultLogger.zl.Sync()
}

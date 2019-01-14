package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger = &logger{
	zaplogger: zap.L(),
	level:     DebugLevel,
}

type logger struct {
	zaplogger *zap.Logger
	level     Level
}

func Error(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= ErrorLevel {
		defaultLogger.zaplogger.Error(msg, fields...)
	}
}

func Errorf(template string, args ...interface{}) {
	if defaultLogger.level >= ErrorLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zaplogger.Error(msg)
	}
}

func Warn(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= WarnLevel {
		defaultLogger.zaplogger.Warn(msg, fields...)
	}
}

func Warnf(template string, args ...interface{}) {
	if defaultLogger.level >= WarnLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zaplogger.Warn(msg)
	}
}

func Info(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= InfoLevel {
		defaultLogger.zaplogger.Info(msg, fields...)
	}
}

func Infof(template string, args ...interface{}) {
	if defaultLogger.level >= InfoLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zaplogger.Info(msg)
	}
}

func Debug(msg string, fields ...zapcore.Field) {
	if defaultLogger.level >= DebugLevel {
		defaultLogger.zaplogger.Debug(msg, fields...)
	}
}

func Debugf(template string, args ...interface{}) {
	if defaultLogger.level >= DebugLevel {
		msg := template
		if len(args) > 0 {
			msg = fmt.Sprintf(template, args...)
		}
		defaultLogger.zaplogger.Debug(msg)
	}
}

func Sync() error {
	return defaultLogger.zaplogger.Sync()
}

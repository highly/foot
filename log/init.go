package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	outputPath      = "stderr"
	errorOutputPath = "stderr"
)

func New(opts ...Option) error {
	o := evaluateOptions(opts)
	errSink, closeErrorSink, err := zap.Open(errorOutputPath)
	if err != nil {
		return err
	}
	outputSink, _, err := zap.Open(outputPath)
	if err != nil {
		closeErrorSink()
		return err
	}
	var rotaterSink zapcore.WriteSyncer
	if o.path != "" {
		rotaterSink = zapcore.AddSync(&lumberjack.Logger{
			Filename:   o.path,
			MaxSize:    o.maxSize,
			MaxBackups: o.maxAge,
			MaxAge:     o.maxBackups,
			LocalTime:  true,
		})
	}
	var sink zapcore.WriteSyncer
	if rotaterSink != nil {
		sink = zapcore.NewMultiWriteSyncer(outputSink, rotaterSink)
	} else {
		sink = outputSink
	}
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     formatDate,
	}
	enc := zapcore.NewJSONEncoder(encCfg)
	core := zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(zap.DebugLevel))
	zapOpts := []zap.Option{
		zap.ErrorOutput(errSink),
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(o.fields) > 0 {
		zapOpts = append(zapOpts, zap.Fields(o.fields...))
	}
	zapLogger := zap.New(core, zapOpts...)
	defaultLogger = &logger{
		zaplogger: zapLogger,
		level:     o.level,
	}
	_ = zap.ReplaceGlobals(zapLogger)
	_ = zap.RedirectStdLog(zapLogger)
	return nil
}

func formatDate(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000
	buf := make([]byte, 27)
	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'
	enc.AppendString(string(buf))
}

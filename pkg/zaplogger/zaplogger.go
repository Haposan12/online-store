package zaplogger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"time"

	"github.com/beego/beego/v2/server/web/context"

	"github.com/bluele/zapslack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	StdFormatLog      = `{app_version:%s, host:%s, path:%s, request_id:%s, request":%s, response":%s}`
	StdFormatErrorLog = `{app_version:%s, host:%s, path:%s, request_id:%s, request:%s, response:%s, error:%s}`
)

type (
	ListErrors struct {
		Error    string
		File     string
		Function string
		Line     int
		Extra    interface{} `json:"extra,omitempty"`
	}
	Fields map[string]interface{}
)

// Logger is our contract for the logger
type Logger interface {
	SetMessageLog(err error, depthList ...int) *ListErrors

	SetMessageErrorToRequestContext(ctx *context.Context, err error, depthList ...int)

	Debug(args ...interface{})

	Debugf(format string, args ...interface{})

	Info(args ...interface{})

	Infof(format string, args ...interface{})

	Warn(args ...interface{})

	Warnf(format string, args ...interface{})

	Error(args ...interface{})

	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})

	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})

	Panicf(format string, args ...interface{})

	With(args ...interface{}) Logger

	WithFields(keyValues Fields) Logger

	Sync() error

	Desugar() *zap.Logger
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger(logPath, slackWebHookUrl string) Logger {

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	})

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(level.CapitalString())
		},
		EncodeTime:     syslogTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, fileSyncer, highPriority),
	)
	logger := zap.New(
		core,
		zap.AddCaller()).Sugar()

	if slackWebHookUrl != "" {
		logger = zap.New(
			core,
			zap.AddCaller(),
			zap.Hooks(zapslack.NewSlackHook(slackWebHookUrl, zap.ErrorLevel).GetHook())).Sugar()
	}

	return &zapLogger{sugaredLogger: logger}
}

func (s zapLogger) SetMessageLog(err error, depthList ...int) *ListErrors {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	le := new(ListErrors)
	if function, file, line, ok := runtime.Caller(depth); ok {
		le.Error = err.Error()
		le.File = file
		le.Function = runtime.FuncForPC(function).Name()
		le.Line = line
	} else {
		le = nil
	}
	return le
}

func (s zapLogger) SetMessageErrorToRequestContext(ctx *context.Context, err error, depthList ...int) {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	if function, file, line, ok := runtime.Caller(depth); ok {
		ctx.Input.SetData("stackTrace", ListErrors{
			Error:    err.Error(),
			File:     file,
			Function: runtime.FuncForPC(function).Name(),
			Line:     line,
		})
	}
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) Sync() error {
	return l.sugaredLogger.Sync()
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

func (l *zapLogger) With(args ...interface{}) Logger {
	newLogger := l.sugaredLogger.With(args...)
	return &zapLogger{newLogger}
}

func (l zapLogger) Desugar() *zap.Logger {
	return l.sugaredLogger.Desugar()
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"synod/conf"
)

const (
	LevelDebug  = zap.DebugLevel
	LevelInfo   = zap.InfoLevel
	LevelWarn   = zap.WarnLevel
	LevelError  = zap.ErrorLevel
	LevelDPanic = zap.DPanicLevel
	LevelPanic  = zap.PanicLevel
	LevelFatal  = zap.FatalLevel
)

type Level = zapcore.Level

type Logger struct {
	desugar *zap.Logger
	sugar   *zap.SugaredLogger
}

func Startup() error {
	logger, err := NewLogger()

	if err != nil {
		return err
	}

	DefaultLogger = logger

	return nil
}

func NewLogger() (*Logger, error) {
	var syncer zapcore.WriteSyncer
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.RFC3339TimeEncoder

	if conf.Bool("app.debug") {
		syncer = zapcore.AddSync(os.Stdout)
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		syncer = zapcore.AddSync(&lumberjack.Logger{
			Filename: conf.String("log.path"),
			MaxSize:  conf.Integer("log.maxSize"),
			MaxAge:   conf.Integer("log.maxAge"),
			Compress: conf.Bool("log.compress"),
		})
		ec.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(ec), syncer, getLogLevelFromConfig())
	zapLogger := zap.New(core, zap.AddCaller())

	logger := &Logger{
		desugar: zapLogger,
		sugar:   zapLogger.Sugar(),
	}

	return logger, nil
}

func getLogLevelFromConfig() Level {
	switch conf.String("log.level") {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "DPanic":
		return LevelDPanic
	case "panic":
		return LevelPanic
	case "fatal":
		return LevelFatal
	}

	return LevelInfo
}

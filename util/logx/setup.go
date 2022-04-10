package logx

import "go.uber.org/zap"

const (
	LevelDebug  = zap.DebugLevel
	LevelInfo   = zap.InfoLevel
	LevelWarn   = zap.WarnLevel
	LevelError  = zap.ErrorLevel
	LevelDPanic = zap.DPanicLevel
	LevelPanic  = zap.PanicLevel
	LevelFatal  = zap.FatalLevel
)

type Logger struct {
	desugar *zap.Logger
	sugar   *zap.SugaredLogger
}

func Setup() error {
	logger, err := NewLogger()

	if err != nil {
		return err
	}

	DefaultLogger = logger

	return nil
}

func NewLogger() (*Logger, error) {
	zapLogger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	logger := &Logger{
		desugar: zapLogger,
		sugar:   zapLogger.Sugar(),
	}

	return logger, nil
}

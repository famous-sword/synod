package logx

func (logger *Logger) Info(msg string, fields ...Field) {
	logger.desugar.Info(msg, fields...)
}

func (logger *Logger) Infow(msg string, fields ...interface{}) {
	logger.sugar.Infow(msg, fields...)
}

func (logger *Logger) Debug(msg string, fields ...Field) {
	logger.desugar.Debug(msg, fields...)
}

func (logger *Logger) Debugw(msg string, fields ...interface{}) {
	logger.sugar.Debugw(msg, fields...)
}

func (logger *Logger) Warn(msg string, fields ...Field) {
	logger.desugar.Warn(msg, fields...)
}

func (logger *Logger) Warnw(msg string, fields ...interface{}) {
	logger.sugar.Warnw(msg, fields...)
}

func (logger *Logger) Error(msg string, fields ...Field) {
	logger.desugar.Error(msg, fields...)
}

func (logger *Logger) Errorw(msg string, fields ...interface{}) {
	logger.sugar.Errorw(msg, fields...)
}

func (logger *Logger) DPanic(msg string, fields ...Field) {
	logger.desugar.DPanic(msg, fields...)
}

func (logger *Logger) DPanicw(msg string, fields ...interface{}) {
	logger.sugar.DPanicw(msg, fields...)
}

func (logger *Logger) Panic(msg string, fields ...Field) {
	logger.desugar.Panic(msg, fields...)
}

func (logger *Logger) Panicw(msg string, fields ...interface{}) {
	logger.sugar.Panicw(msg, fields...)
}

func (logger *Logger) Fatal(msg string, fields ...Field) {
	logger.desugar.Fatal(msg, fields...)
}

func (logger *Logger) Fatalw(msg string, fields ...interface{}) {
	logger.sugar.Fatalw(msg, fields...)
}

func (logger *Logger) Flush() {
	_ = logger.desugar.Sync()
	_ = logger.sugar.Sync()
}

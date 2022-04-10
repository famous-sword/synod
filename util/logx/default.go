package logx

var DefaultLogger *Logger

func Info(msg string, fields ...Field) {
	DefaultLogger.Info(msg, fields...)
}

func Infow(msg string, fields ...interface{}) {
	DefaultLogger.Infow(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	DefaultLogger.Debug(msg, fields...)
}

func Debugw(msg string, fields ...interface{}) {
	DefaultLogger.Debugw(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	DefaultLogger.Warn(msg, fields...)
}

func Warnw(msg string, fields ...interface{}) {
	DefaultLogger.Warnw(msg, fields...)
}

func Error(msg string, fields ...Field) {
	DefaultLogger.Error(msg, fields...)
}

func Errorw(msg string, fields ...interface{}) {
	DefaultLogger.Errorw(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	DefaultLogger.DPanic(msg, fields...)
}

func DPanicw(msg string, fields ...interface{}) {
	DefaultLogger.DPanicw(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	DefaultLogger.Panic(msg, fields...)
}

func Panicw(msg string, fields ...interface{}) {
	DefaultLogger.Panicw(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	DefaultLogger.Fatal(msg, fields...)
}

func Fatalw(msg string, fields ...interface{}) {
	DefaultLogger.Fatalw(msg, fields...)
}

func Flush() {
	DefaultLogger.Flush()
}

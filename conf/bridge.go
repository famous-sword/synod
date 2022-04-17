package conf

func String(key string) string {
	return repository.String(key)
}

func StringSlice(key string) []string {
	return repository.StringSlice(key)
}

func IntSlice(key string) []int {
	return repository.IntSlice(key)
}

func Integer(key string) int {
	return repository.Integer(key)
}

func Bool(key string) bool {
	return repository.Bool(key)
}

func Int32(key string) int32 {
	return repository.Int32(key)
}

func Int64(key string) int64 {
	return repository.Int64(key)
}

func Uint(key string) uint {
	return repository.Uint(key)
}

func Uint32(key string) uint32 {
	return repository.Uint32(key)
}

func Uint64(key string) uint64 {
	return repository.Uint64(key)
}

func Float(key string) float64 {
	return repository.Float(key)
}

func Set(key string, value interface{}) {
	repository.Set(key, value)
}

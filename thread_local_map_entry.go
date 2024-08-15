package routine

type entry any

func entryValue[T any](e entry) T {
	if e == nil {
		var defaultValue T
		return defaultValue
	}
	return e.(T)
}

func entryAssert[T any](e entry) (T, bool) {
	v, ok := e.(T)
	return v, ok
}

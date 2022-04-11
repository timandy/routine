package routine

// Goid return the current goroutine's unique id.
func Goid() int64 {
	return getg().goid
}

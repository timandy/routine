package routine

// Goid return the current goroutine's unique id.
func Goid() uint64 {
	return getg().goid()
}

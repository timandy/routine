package routine

// ThreadLocal provides goroutine-local variables.
type ThreadLocal interface {
	// Get returns the value in the current goroutine's local threadLocals or inheritableThreadLocals, if it was set before.
	Get() any

	// Set copy the value into the current goroutine's local threadLocals or inheritableThreadLocals.
	Set(value any)

	// Remove delete the value from the current goroutine's local threadLocals or inheritableThreadLocals.
	Remove()
}

// Supplier provides a function that returns a value of type any.
type Supplier func() any

// NewThreadLocal create and return a new ThreadLocal instance.
// The initial value is nil.
func NewThreadLocal() ThreadLocal {
	return &threadLocal{index: nextThreadLocalIndex()}
}

// NewThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value is determined by invoking the supplier method.
func NewThreadLocalWithInitial(supplier Supplier) ThreadLocal {
	return &threadLocal{index: nextThreadLocalIndex(), supplier: supplier}
}

// NewInheritableThreadLocal create and return a new ThreadLocal instance.
// The initial value is nil.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
func NewInheritableThreadLocal() ThreadLocal {
	return &inheritableThreadLocal{index: nextInheritableThreadLocalIndex()}
}

// NewInheritableThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value is determined by invoking the supplier method.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
func NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal {
	return &inheritableThreadLocal{index: nextInheritableThreadLocalIndex(), supplier: supplier}
}

package routine

// ThreadLocal provides goroutine-local variables.
type ThreadLocal interface {
	// Get returns the value in the current goroutine's local threadLocals or inheritableThreadLocals, if it was set before.
	Get() Any

	// Set copy the value into the current goroutine's local threadLocals or inheritableThreadLocals.
	Set(value Any)

	// Remove delete the value from the current goroutine's local threadLocals or inheritableThreadLocals.
	Remove()
}

// Supplier provides a function which return Any type result.
type Supplier = func() Any

// NewThreadLocal create and return a new ThreadLocal instance.
// The initial value is nil.
func NewThreadLocal() ThreadLocal {
	return &threadLocal{id: nextThreadLocalId()}
}

// NewThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value is determined by invoking the supplier method.
func NewThreadLocalWithInitial(supplier Supplier) ThreadLocal {
	return &threadLocal{id: nextThreadLocalId(), supplier: supplier}
}

// NewInheritableThreadLocal create and return a new ThreadLocal instance.
// The initial value is nil.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
func NewInheritableThreadLocal() ThreadLocal {
	return &inheritableThreadLocal{id: nextInheritableThreadLocalId()}
}

// NewInheritableThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value is determined by invoking the supplier method.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
func NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal {
	return &inheritableThreadLocal{id: nextInheritableThreadLocalId(), supplier: supplier}
}

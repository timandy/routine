package routine

// ThreadLocal provides goroutine-local variables.
type ThreadLocal[T any] interface {
	// Get returns the value in the current goroutine's local threadLocals or inheritableThreadLocals, if it was set before.
	Get() T

	// Set copy the value into the current goroutine's local threadLocals or inheritableThreadLocals.
	Set(value T)

	// Remove delete the value from the current goroutine's local threadLocals or inheritableThreadLocals.
	Remove()
}

// Supplier provides a function that returns a value of type T.
type Supplier[T any] func() T

// NewThreadLocal create and return a new ThreadLocal instance.
// The initial value stored with the default value of type T.
func NewThreadLocal[T any]() ThreadLocal[T] {
	return &threadLocal[T]{index: nextThreadLocalIndex()}
}

// NewThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value stored as the return value of the method supplier.
func NewThreadLocalWithInitial[T any](supplier Supplier[T]) ThreadLocal[T] {
	return &threadLocal[T]{index: nextThreadLocalIndex(), supplier: supplier}
}

// NewInheritableThreadLocal create and return a new ThreadLocal instance.
// The initial value stored with the default value of type T.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
// The value can be captured to FutureTask which created by WrapTask, WrapWaitTask, WrapWaitResultTask methods.
func NewInheritableThreadLocal[T any]() ThreadLocal[T] {
	return &inheritableThreadLocal[T]{index: nextInheritableThreadLocalIndex()}
}

// NewInheritableThreadLocalWithInitial create and return a new ThreadLocal instance.
// The initial value stored as the return value of the method supplier.
// The value can be inherited to sub goroutines witch started by Go, GoWait, GoWaitResult methods.
// The value can be captured to FutureTask which created by WrapTask, WrapWaitTask, WrapWaitResultTask methods.
func NewInheritableThreadLocalWithInitial[T any](supplier Supplier[T]) ThreadLocal[T] {
	return &inheritableThreadLocal[T]{index: nextInheritableThreadLocalIndex(), supplier: supplier}
}

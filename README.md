# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://codecov.io/gh/timandy/routine)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/timandy/routine)

> [中文版](README_zh.md)

`routine` encapsulates and provides some easy-to-use, high-performance `goroutine` context access interfaces, which can help you access coroutine context information more elegantly, but you may also open Pandora's Box.

# Introduce

The `Golang` language has been sparing no effort to shield developers from the concept of coroutine context from the beginning of its design, including the acquisition of coroutine `goid`, the state of the coroutine within the process, and the storage of coroutine context.

If you have used other languages such as `C++/Java/...`, then you must be familiar with `ThreadLocal`, and after starting to use `Golang`, you will definitely feel confused and distressed by the lack of convenient functions similar to `ThreadLocal`. Of course, you can choose to use `Context`, let it carry all the context information, appear in the first input parameter of all functions, and then shuttle around in your system.

The core goal of `routine` is to open up another path: to introduce `goroutine local storage` into the world of `Golang`, and at the same time expose the coroutine information to meet the needs of some people.

# Usage & Demo

This chapter briefly introduces how to install and use the `routine` library.

## Install

```bash
go get github.com/timandy/routine
```

## Use `goid`

The following code simply demonstrates the use of `routine.Goid()` and `routine.AllGoids()`:

```go
package main

import (
	"fmt"
	"github.com/timandy/routine"
	"time"
)

func main() {
	go func() {
		time.Sleep(time.Second)
	}()
	goid := routine.Goid()
	goids := routine.AllGoids()
	fmt.Printf("curr goid: %d\n", goid)
	fmt.Printf("all goids: %v\n", goids)
}
```

In this example, the `main` function starts a new coroutine, so `Goid()` returns the main coroutine `1`, and `AllGoids()` returns the main coroutine and coroutine `18`:

```text
curr goid: 1
all goids: [1 18]
```

## Use `ThreadLocal`

The following code briefly demonstrates `ThreadLocal`'s creation, setting, getting, spreading across coroutines, etc.:

```go
package main

import (
	"fmt"
	"github.com/timandy/routine"
	"time"
)

var threadLocal = routine.NewThreadLocal()
var inheritableThreadLocal = routine.NewInheritableThreadLocal()

func main() {
	threadLocal.Set("hello world")
	inheritableThreadLocal.Set("Hello world2")
	fmt.Println("threadLocal:", threadLocal.Get())
	fmt.Println("inheritableThreadLocal:", inheritableThreadLocal.Get())

	// Other goroutines cannot read the previously assigned "hello world"
	go func() {
		fmt.Println("threadLocal in goroutine:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine:", inheritableThreadLocal.Get())
	}()

	// However, a new goroutine can be started via the Go function. All the inheritable variables of the current main goroutine can be passed away automatically.
	routine.Go(func() {
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

	time.Sleep(time.Second)
}
```

The execution result is:

```text
threadLocal: hello world
inheritableThreadLocal: Hello world2
threadLocal in goroutine: <nil>
inheritableThreadLocal in goroutine: <nil>
threadLocal in goroutine by Go: <nil>
inheritableThreadLocal in goroutine by Go: Hello world2
```

# API

This chapter introduces in detail all the interfaces encapsulated by the `routine` library, as well as their core functions and implementation methods.

## `Goid() int64`

Get the `goid` of the current `goroutine`.

Under normal circumstances, `Goid()` first tries to obtain it directly through `go_tls`. This operation has extremely high performance and the time-consuming is usually only one-fifth of `rand.Int()`.

If an error such as version incompatibility occurs, `Goid()` will try to downgrade, that is, parse it from the `runtime.Stack` information. At this time, the performance will drop sharply by about a thousand times, but it can ensure that the function is normally available.

## `AllGoids() []int64`

Get the `goid` of all active `goroutine` of the current process.

In `go 1.15` and older versions, `AllGoids()` will try to parse and get all the coroutine information from the `runtime.Stack` information, but this operation is very inefficient, and it is not recommended using it in high-frequency logic. .

In versions after `go 1.16`, `AllGoids()` will directly read the global coroutine pool information of `runtime` through `native`, which has greatly improved performance, but considering the production environment There may be tens of thousands or millions of coroutines, so it is still not recommended using it at high frequencies.

## `NewThreadLocal() ThreadLocal`

Creates a new `ThreadLocal` instance with a stored default value of `nil`.

## `NewThreadLocalWithInitial(supplier Supplier) ThreadLocal`

Creates a new `ThreadLocal` instance with default values stored by calling `supplier()`.

## `NewInheritableThreadLocal() ThreadLocal`

Creates a new `ThreadLocal` instance with a stored default value of `nil`. When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine .

## `NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal`

Creates a new `ThreadLocal` instance with stored default values generated by calling `supplier()`. When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine .

## `Go(fun func())`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. Any `panic` while the child coroutine is executing will be caught and the stack automatically printed.

## `GoWait(fun func()) Feature`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. You can wait for the sub-coroutine to finish executing through the `Feature.Get()` method that returns a value. Any `panic` while the child coroutine is executing will be caught and thrown again when `Feature.Get()` is called.

## `GoWaitResult(fun func() Any) Feature`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. You can wait for the sub-coroutine to finish executing and get the return value through the `Feature.Get()` method of the return value. Any `panic` while the child coroutine is executing will be caught and thrown again when `Feature.Get()` is called.

[More API Documentation](https://pkg.go.dev/github.com/timandy/routine#section-documentation)

# Garbage Collection

The `routine` library internally maintains the global `globalMap` variable, which stores all the context variable information of the coroutine, and performs variable addressing mapping based on the `goid` of the coroutine and the `ptr` of the coroutine variable when reading and writing.

In the entire life cycle of a process, it may be created by destroying countless coroutines, so how to clean up the context variables of these `dead` coroutines?

To solve this problem, a global `gcTimer` is allocated internally by `routine`. This timer will be started when `globalMap` needs to be cleaned up. It scans and cleans up the context variables cached by `dead` coroutine in `globalMap` at regular intervals, to avoid possible hidden dangers of memory leaks.

# License

MIT

# Thanks

This lib is forked from [go-eden/routine](https://github.com/go-eden/routine). Thank go-eden for his great work!

# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://app.codecov.io/gh/timandy/routine)
[![Documentation](https://pkg.go.dev/badge/github.com/timandy/routine.svg)](https://pkg.go.dev/github.com/timandy/routine)
[![Release](https://img.shields.io/github/release/timandy/routine.svg)](https://github.com/timandy/routine/releases)
[![License](https://img.shields.io/github/license/timandy/routine.svg)](https://github.com/timandy/routine/blob/main/LICENSE)

> [中文版](README_zh.md)

`routine` encapsulates and provides some easy-to-use, non-competitive, high-performance `goroutine` context access interfaces, which can help you access coroutine context information more gracefully.

# Introduce

From the very beginning of its design, the `Golang` language has spared no effort to shield the concept of coroutine context from developers, including the acquisition of coroutine `goid`, the state of coroutine within the process, and the storage of coroutine context.

If you have used other languages such as `C++`, `Java` and so on, then you must be familiar with `ThreadLocal`, but after starting to use `Golang`, you will be deeply confused and distressed by the lack of convenient functions like `ThreadLocal`.

Of course, you can choose to use `Context`, which carries all the context information, appears in the first input parameter of all functions, and then shuttles around your system.

And the core goal of `routine` is to open up another way: Introduce `goroutine local storage` to the `Golang` world.

# Usage & Demo

This chapter briefly introduces how to install and use the `routine` library.

## Install

```bash
go get github.com/timandy/routine
```

## Use `goid`

The following code simply demonstrates the use of `routine.Goid()`:

```go
package main

import (
	"fmt"
	"github.com/timandy/routine"
	"time"
)

func main() {
	goid := routine.Goid()
	fmt.Printf("cur goid: %v\n", goid)
	go func() {
		goid := routine.Goid()
		fmt.Printf("sub goid: %v\n", goid)
	}()

	// Wait for the sub-coroutine to finish executing.
	time.Sleep(time.Second)
}
```

In this example, the `main` function starts a new coroutine, so `Goid()` returns the main coroutine `1` and the child coroutine `6`:

```text
cur goid: 1
sub goid: 6
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

	// The child coroutine cannot read the previously assigned "hello world".
	go func() {
		fmt.Println("threadLocal in goroutine:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine:", inheritableThreadLocal.Get())
	}()

	// However, a new sub-coroutine can be started via the Go/GoWait/GoWaitResul function, and all inheritable variables of the current coroutine can be passed automatically.
	routine.Go(func() {
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

	// Wait for the sub-coroutine to finish executing.
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

It can be obtained directly through assembly code under `386`, `amd64`, `armv6`, `armv7`, `arm64` architectures. This operation has extremely high performance and the time-consuming is usually only one-fifth of `rand.Int()`.

## `NewThreadLocal() ThreadLocal`

Creates a new `ThreadLocal` instance with a stored default value of `nil`.

## `NewThreadLocalWithInitial(supplier Supplier) ThreadLocal`

Creates a new `ThreadLocal` instance with default values stored by calling `supplier()`.

## `NewInheritableThreadLocal() ThreadLocal`

Creates a new `ThreadLocal` instance with a stored default value of `nil`. When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine.

## `NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal`

Creates a new `ThreadLocal` instance with stored default values generated by calling `supplier()`. When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine.

## `Go(fun func())`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. Any `panic` while the child coroutine is executing will be caught and the stack automatically printed.

## `GoWait(fun func()) Feature`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. You can wait for the sub-coroutine to finish executing through the `Feature.Get()` method that returns a value. Any `panic` while the child coroutine is executing will be caught and thrown again when `Feature.Get()` is called.

## `GoWaitResult(fun func() Any) Feature`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine. You can wait for the sub-coroutine to finish executing and get the return value through the `Feature.Get()` method of the return value. Any `panic` while the child coroutine is executing will be caught and thrown again when `Feature.Get()` is called.

[More API Documentation](https://pkg.go.dev/github.com/timandy/routine#section-documentation)

# Garbage Collection

`routine` allocates a `thread` structure for each coroutine, which stores context variable information related to the coroutine.

A pointer to this structure is stored on the `g.labels` field of the coroutine structure.

When the coroutine finishes executing and exits, `g.labels` will be set to `nil`, no longer referencing the `thread` structure.

The `thread` structure will be collected at the next `GC`.

If the data stored in `thread` is not additionally referenced, these data will be collected together.

# Support Grid

|             | **`darwin`** | **`linux`** | **`windows`** |             |
|------------:|:------------:|:-----------:|:-------------:|:------------|
|   **`386`** |              |      ✅      |       ✅       | **`386`**   |
| **`amd64`** |      ✅       |      ✅      |       ✅       | **`amd64`** |
| **`armv6`** |              |      ✅      |               | **`armv6`** |
| **`armv7`** |              |      ✅      |               | **`armv7`** |
| **`arm64`** |              |      ✅      |               | **`arm64`** |
|             | **`darwin`** | **`linux`** | **`windows`** |             |

✅: Supported

# Thanks

`routine` is forked from [go-eden/routine](https://github.com/go-eden/routine), thanks to the original author for his contribution!

# *License*

`routine` is released under the [Apache License 2.0](LICENSE).

```
Copyright 2021-2022 TimAndy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

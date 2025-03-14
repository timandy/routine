# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://app.codecov.io/gh/timandy/routine)
[![Go Report Card](https://goreportcard.com/badge/github.com/timandy/routine)](https://goreportcard.com/report/github.com/timandy/routine)
[![Documentation](https://pkg.go.dev/badge/github.com/timandy/routine.svg)](https://pkg.go.dev/github.com/timandy/routine)
[![Release](https://img.shields.io/github/release/timandy/routine.svg)](https://github.com/timandy/routine/releases)
[![License](https://img.shields.io/github/license/timandy/routine.svg)](https://github.com/timandy/routine/blob/main/LICENSE)

> [中文版](README_zh.md)

`routine` encapsulates and provides some easy-to-use, non-competitive, high-performance `goroutine` context access interfaces, which can help you access coroutine context information more gracefully.

# :house:Introduce

From the very beginning of its design, the `Golang` language has spared no effort to shield the concept of coroutine context from developers, including the acquisition of coroutine `goid`, the state of coroutine within the process, and the storage of coroutine context.

If you have used other languages such as `C++`, `Java` and so on, then you must be familiar with `ThreadLocal`, but after starting to use `Golang`, you will be deeply confused and distressed by the lack of convenient functions like `ThreadLocal`.

Of course, you can choose to use `Context`, which carries all the context information, appears in the first input parameter of all functions, and then shuttles around your system.

And the core goal of `routine` is to open up another way: Introduce `goroutine local storage` to the `Golang` world.

# :loudspeaker:Update Notice

:fire:**Version `1.1.5` introduces a new static mode.**

- :rocket:Performance improved by over `20%`.

- :rocket:Memory access is now safer.

- :exclamation:The compile command requires additional parameters `-a -toolexec='routinex -v'`.

For more details, visit: [RoutineX Compiler](https://github.com/timandy/routinex)

# :hammer_and_wrench:Usage & Demo

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
	"time"

	"github.com/timandy/routine"
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
	"time"

	"github.com/timandy/routine"
)

var threadLocal = routine.NewThreadLocal[string]()
var inheritableThreadLocal = routine.NewInheritableThreadLocal[string]()

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

	// However, a new sub-coroutine can be started via the Go/GoWait/GoWaitResult function, and all inheritable variables of the current coroutine can be passed automatically.
	routine.Go(func() {
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

	// You can also create a task via the WrapTask/WrapWaitTask/WrapWaitResultTask function, and all inheritable variables of the current coroutine can be automatically captured.
	task := routine.WrapTask(func() {
		fmt.Println("threadLocal in task by WrapTask:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in task by WrapTask:", inheritableThreadLocal.Get())
	})
	go task.Run()

	// Wait for the sub-coroutine to finish executing.
	time.Sleep(time.Second)
}
```

The execution result is:

```text
threadLocal: hello world
inheritableThreadLocal: Hello world2
threadLocal in goroutine:
inheritableThreadLocal in goroutine:
threadLocal in goroutine by Go:
inheritableThreadLocal in goroutine by Go: Hello world2
threadLocal in task by WrapTask:
inheritableThreadLocal in task by WrapTask: Hello world2
```

# :books:API

This chapter introduces in detail all the interfaces encapsulated by the `routine` library, as well as their core functions and implementation methods.

## `Goid() uint64`

Get the `goid` of the current `goroutine`.

It can be obtained directly through assembly code under `386`, `amd64`, `armv6`, `armv7`, `arm64`, `loong64`, `mips`, `mipsle`, `mips64`, `mips64le`, `ppc64`, `ppc64le`, `riscv64`, `s390x`, `wasm` architectures. This operation has extremely high performance and the time-consuming is usually only one-fifth of `rand.Int()`.

## `NewThreadLocal[T any]() ThreadLocal[T]`

Create a new `ThreadLocal[T]` instance with the initial value stored with the default value of type `T`.

## `NewThreadLocalWithInitial[T any](supplier Supplier[T]) ThreadLocal[T]`

Create a new `ThreadLocal[T]` instance with the initial value stored as the return value of the method `supplier()`.

## `NewInheritableThreadLocal[T any]() ThreadLocal[T]`

Create a new `ThreadLocal[T]` instance with the initial value stored with the default value of type `T`.
When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine.
When a new task is created via `WrapTask()`, `WrapWaitTask()` or `WrapWaitResultTask()`, the value of the current coroutine is captured to the new task.

## `NewInheritableThreadLocalWithInitial[T any](supplier Supplier[T]) ThreadLocal[T]`

Create a new `ThreadLocal[T]` instance with the initial value stored as the return value of the method `supplier()`.
When a new coroutine is started via `Go()`, `GoWait()` or `GoWaitResult()`, the value of the current coroutine is copied to the new coroutine.
When a new task is created via `WrapTask()`, `WrapWaitTask()` or `WrapWaitResultTask()`, the value of the current coroutine is captured to the new task.

## `WrapTask(fun Runnable) FutureTask[any]`

Create a new task and capture the `inheritableThreadLocals` from the current goroutine.
This function returns a `FutureTask` instance, but the return task will not run automatically.
You can run it in a sub-goroutine or goroutine-pool by `FutureTask.Run()` method, wait by `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.
When the returned task run `panic` will be caught and error stack will be printed, the `panic` will be trigger again when calling `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.

## `WrapWaitTask(fun CancelRunnable) FutureTask[any]`

Create a new task and capture the `inheritableThreadLocals` from the current goroutine.
This function returns a `FutureTask` instance, but the return task will not run automatically.
You can run it in a sub-goroutine or goroutine-pool by `FutureTask.Run()` method, wait by `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.
When the returned task run `panic` will be caught, the `panic` will be trigger again when calling `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.

## `WrapWaitResultTask[TResult any](fun CancelCallable[TResult]) FutureTask[TResult]`

Create a new task and capture the `inheritableThreadLocals` from the current goroutine.
This function returns a `FutureTask` instance, but the return task will not run automatically.
You can run it in a sub-goroutine or goroutine-pool by `FutureTask.Run()` method, wait and get result by `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.
When the returned task run `panic` will be caught, the `panic` will be trigger again when calling `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method.

## `Go(fun Runnable)`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine.
Any `panic` while the child coroutine is executing will be caught and the stack automatically printed.

## `GoWait(fun CancelRunnable) FutureTask[any]`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine.
You can wait for the sub-coroutine to finish executing through the `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method that returns a value.
Any `panic` while the child coroutine is executing will be caught and thrown again when `FutureTask.Get()` or `FutureTask.GetWithTimeout()` is called.

## `GoWaitResult[TResult any](fun CancelCallable[TResult]) FutureTask[TResult]`

Start a new coroutine and automatically copy all contextual `inheritableThreadLocals` data of the current coroutine to the new coroutine.
You can wait for the sub-coroutine to finish executing and get the return value through the `FutureTask.Get()` or `FutureTask.GetWithTimeout()` method of the return value.
Any `panic` while the child coroutine is executing will be caught and thrown again when `FutureTask.Get()` or `FutureTask.GetWithTimeout()` is called.

[More API Documentation](https://pkg.go.dev/github.com/timandy/routine#section-documentation)

# :wastebasket:Garbage Collection

`routine` allocates a `thread` structure for each coroutine, which stores context variable information related to the coroutine.

A pointer to this structure is stored on the `g.labels` field of the coroutine structure.

When the coroutine finishes executing and exits, `g.labels` will be set to `nil`, no longer referencing the `thread` structure.

The `thread` structure will be collected at the next `GC`.

If the data stored in `thread` is not additionally referenced, these data will be collected together.

# :globe_with_meridians:Support Grid

|                | **`darwin`** | **`linux`** | **`windows`** | **`freebsd`** | **`js`** |                |
|---------------:|:------------:|:-----------:|:-------------:|:-------------:|:--------:|:---------------|
|      **`386`** |              |      ✅      |       ✅       |       ✅       |          | **`386`**      |
|    **`amd64`** |      ✅       |      ✅      |       ✅       |       ✅       |          | **`amd64`**    |
|    **`armv6`** |              |      ✅      |               |               |          | **`armv6`**    |
|    **`armv7`** |              |      ✅      |               |               |          | **`armv7`**    |
|    **`arm64`** |      ✅       |      ✅      |               |               |          | **`arm64`**    |
|  **`loong64`** |              |      ✅      |               |               |          | **`loong64`**  |
|     **`mips`** |              |      ✅      |               |               |          | **`mips`**     |
|   **`mipsle`** |              |      ✅      |               |               |          | **`mipsle`**   |
|   **`mips64`** |              |      ✅      |               |               |          | **`mips64`**   |
| **`mips64le`** |              |      ✅      |               |               |          | **`mips64le`** |
|    **`ppc64`** |              |      ✅      |               |               |          | **`ppc64`**    |
|  **`ppc64le`** |              |      ✅      |               |               |          | **`ppc64le`**  |
|  **`riscv64`** |              |      ✅      |               |               |          | **`riscv64`**  |
|    **`s390x`** |              |      ✅      |               |               |          | **`s390x`**    |
|     **`wasm`** |              |             |               |               |    ✅     | **`wasm`**     |
|                | **`darwin`** | **`linux`** | **`windows`** | **`freebsd`** | **`js`** |                |

✅: Supported

# :pray:Thanks

Thanks to all [contributors](https://github.com/timandy/routine/graphs/contributors) for their contributions!

# :scroll:*License*

`routine` is released under the [Apache License 2.0](LICENSE).

```
Copyright 2021-2025 TimAndy

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

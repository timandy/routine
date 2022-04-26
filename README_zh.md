# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://app.codecov.io/gh/timandy/routine)
[![Documentation](https://pkg.go.dev/badge/github.com/timandy/routine.svg)](https://pkg.go.dev/github.com/timandy/routine)
[![Release](https://img.shields.io/github/release/timandy/routine.svg)](https://github.com/timandy/routine/releases)
[![License](https://img.shields.io/github/license/timandy/routine.svg)](https://github.com/timandy/routine/blob/main/LICENSE)

> [English Version](README.md)

`routine`封装并提供了一些易用、无竞争、高性能的`goroutine`上下文访问接口，它可以帮助你更优雅地访问协程上下文信息。

# 介绍

`Golang`语言从设计之初，就一直在不遗余力地向开发者屏蔽协程上下文的概念，包括协程`goid`的获取、进程内部协程状态、协程上下文存储等。

如果你使用过其他语言如`C++`、`Java`等，那么你一定很熟悉`ThreadLocal`，而在开始使用`Golang`之后，你一定会为缺少类似`ThreadLocal`的便捷功能而深感困惑与苦恼。

当然你可以选择使用`Context`，让它携带着全部上下文信息，在所有函数的第一个输入参数中出现，然后在你的系统中到处穿梭。

而`routine`的核心目标就是开辟另一条路：将`goroutine local storage`引入`Golang`世界。

# 使用演示

此章节简要介绍如何安装与使用`routine`库。

## 安装

```bash
go get github.com/timandy/routine
```

## 使用`goid`

以下代码简单演示了`routine.Goid()`的使用：

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

	// 等待子协程执行完。
	time.Sleep(time.Second)
}
```

此例中`main`函数启动了一个新的协程，因此`Goid()`返回了主协程`1`和子协程`6`:

```text
cur goid: 1
sub goid: 6
```

## 使用`ThreadLocal`

以下代码简单演示了`ThreadLocal`的创建、设置、获取、跨协程传播等：

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

	// 子协程无法读取之前赋值的“hello world”。
	go func() {
		fmt.Println("threadLocal in goroutine:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine:", inheritableThreadLocal.Get())
	}()

	// 但是，可以通过 Go/GoWait/GoWaitResul 函数启动一个新的子协程，当前协程的所有可继承变量都可以自动传递。
	routine.Go(func() {
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

	// 等待子协程执行完。
	time.Sleep(time.Second)
}
```

执行结果为：

```text
threadLocal: hello world
inheritableThreadLocal: Hello world2
threadLocal in goroutine: <nil>
inheritableThreadLocal in goroutine: <nil>
threadLocal in goroutine by Go: <nil>
inheritableThreadLocal in goroutine by Go: Hello world2
```

# API文档

此章节详细介绍了`routine`库封装的全部接口，以及它们的核心功能、实现方式等。

## `Goid() int64`

获取当前`goroutine`的`goid`。

在`386`、`amd64`、`armv6`、`armv7`、`arm64`架构下通过汇编代码直接获取，此操作性能极高，耗时通常只相当于`rand.Int()`的五分之一。

## `NewThreadLocal() ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值为`nil`。

## `NewThreadLocalWithInitial(supplier Supplier) ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值会通过调用`supplier()`生成。

## `NewInheritableThreadLocal() ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值为`nil`。当通过`Go()`、`GoWait()`或`GoWaitResult()`启动新协程时，当前协程的值会被复制到新协程。

## `NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值会通过调用`supplier()`生成。当通过`Go()`、`GoWait()`或`GoWaitResult()`启动新协程时，当前协程的值会被复制到新协程。

## `Go(fun func())`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。子协程执行时的任何`panic`都会被捕获并自动打印堆栈。

## `GoWait(fun func()) Feature`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。可以通过返回值的`Feature.Get()`方法等待子协程执行完毕。子协程执行时的任何`panic`都会被捕获并在调用`Feature.Get()`时再次抛出。

## `GoWaitResult(fun func() Any) Feature`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。可以通过返回值的`Feature.Get()`方法等待子协程执行完毕并获取返回值。子协程执行时的任何`panic`都会被捕获并在调用`Feature.Get()`时再次抛出。

[更多API文档](https://pkg.go.dev/github.com/timandy/routine#section-documentation)

# 垃圾回收

`routine`为每个协程分配了一个`thread`结构，它存储了协程相关的上下文变量信息。

指向该结构的指针存储在协程结构的`g.labels`字段上。

当协程执行完毕退出时，`g.labels`将被设置为`nil`，不再引用`thread`结构。

`thread`结构将在下次`GC`时被回收。

如果`thread`中存储的数据也没有额外被引用，这些数据将被一并回收。

# 支持网格

|             | **`darwin`** | **`linux`** | **`windows`** |             |
|------------:|:------------:|:-----------:|:-------------:|:------------|
|   **`386`** |              |      ✅      |       ✅       | **`386`**   |
| **`amd64`** |      ✅       |      ✅      |       ✅       | **`amd64`** |
| **`armv6`** |              |      ✅      |               | **`armv6`** |
| **`armv7`** |              |      ✅      |               | **`armv7`** |
| **`arm64`** |              |      ✅      |               | **`arm64`** |
|             | **`darwin`** | **`linux`** | **`windows`** |             |

✅：支持

# 鸣谢

`routine`是从 [go-eden/routine](https://github.com/go-eden/routine) 分支出来的，感谢原作者的贡献！

# *许可证*

`routine`是在 [Apache License 2.0](LICENSE) 下发布的。

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

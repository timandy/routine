# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://codecov.io/gh/timandy/routine)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/timandy/routine)

> [English Version](README.md)

`routine`封装并提供了一些易用、高性能的`goroutine`上下文访问接口，它可以帮助你更优雅地访问协程上下文信息，但你也可能就此打开了潘多拉魔盒。

# 介绍

`Golang`语言从设计之初，就一直在不遗余力地向开发者屏蔽协程上下文的概念，包括协程`goid`的获取、进程内部协程状态、协程上下文存储等。

如果你使用过其他语言如`C++/Java`等，那么你一定很熟悉`ThreadLocal`，而在开始使用`Golang`之后，你一定会为缺少类似`ThreadLocal`的便捷功能而深感困惑与苦恼。 当然你可以选择使用`Context`，让它携带着全部上下文信息，在所有函数的第一个输入参数中出现，然后在你的系统中到处穿梭。

而`routine`的核心目标就是开辟另一条路：将`goroutine local storage`引入`Golang`世界，同时也将协程信息暴露出来，以满足某些人可能有的需求。

# 使用演示

此章节简要介绍如何安装与使用`routine`库。

## 安装

```bash
go get github.com/timandy/routine
```

## 使用`goid`

以下代码简单演示了`routine.Goid()`与`routine.AllGoids()`的使用：

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
	fmt.Printf("curr goid: %v\n", goid)
	fmt.Printf("all goids: %v\n", goids)
	fmt.Print("each goid:")
	routine.ForeachGoid(func(goid int64) {
		fmt.Printf(" %v", goid)
	})
}
```

此例中`main`函数启动了一个新的协程，因此`Goid()`返回了主协程`1`，`AllGoids()`返回了主协程及协程`18`，`ForeachGoid()`依次返回了主协程及协程`18`:

```text
curr goid: 1
all goids: [1 18]
each goid: 1 18
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

	// 其他协程无法读取之前赋值的“hello world”。
	go func() {
		fmt.Println("threadLocal in goroutine:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine:", inheritableThreadLocal.Get())
	}()

	// 但是，可以通过 Go 函数启动一个新的 goroutine。当前主 goroutine 的所有可继承变量都可以自动传递。
	routine.Go(func() {
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

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

在正常情况下，`Goid()`优先尝试通过`go_tls`的方式直接获取，此操作性能极高，耗时通常只相当于`rand.Int()`的五分之一。

若出现版本不兼容等错误时，`Goid()`会尝试降级，即从`runtime.Stack`信息中解析获取，此时性能会急剧下降约千倍，但它可以保证功能正常可用。

## `AllGoids() []int64`

获取当前进程全部活跃`goroutine`的`goid`。在执行过程中添加新的`goid`可能会被遗漏。

在`go 1.12`及更旧的版本中，`AllGoids()`会尝试从`runtime.Stack`信息中解析获取全部协程信息，但此操作非常低效，非常不建议在高频逻辑中使用。

在`go 1.13`之后的版本中，`AllGoids()`会通过`native`的方式直接读取`runtime`的全局协程池信息，在性能上得到了极大的提高，但考虑到生产环境中可能有万、百万级的协程数量，因此仍不建议在高频使用它。

## `ForeachGoid(fun func(goid int64))`

为当前进程全部活跃`goroutine`的`goid`执行指定函数。在执行过程中添加新的`goid`可能会被遗漏。

获取`goid`的方式同`AllGoids() []int64`。

## `NewThreadLocal() ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值为`nil`。

## `NewThreadLocalWithInitial(supplier Supplier) ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值会通过调用`supplier()`生成。

## `NewInheritableThreadLocal() ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值为`nil`。当通过 `Go()`、`GoWait()`或`GoWaitResult()` 启动新协程时，当前协程的值会被复制到新协程。

## `NewInheritableThreadLocalWithInitial(supplier Supplier) ThreadLocal`

创建一个新的`ThreadLocal`实例，其存储的默认值会通过调用`supplier()`生成。当通过 `Go()`、`GoWait()`或`GoWaitResult()` 启动新协程时，当前协程的值会被复制到新协程。

## `Go(fun func())`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。子协程执行时的任何`panic`都会被捕获并自动打印堆栈。

## `GoWait(fun func()) Feature`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。可以通过返回值的`Feature.Get()`方法等待子协程执行完毕。子协程执行时的任何`panic`都会被捕获并在调用`Feature.Get()`时再次抛出。

## `GoWaitResult(fun func() Any) Feature`

启动一个新的协程，同时自动将当前协程的全部上下文`inheritableThreadLocals`数据复制至新协程。可以通过返回值的`Feature.Get()`方法等待子协程执行完毕并获取返回值。子协程执行时的任何`panic`都会被捕获并在调用`Feature.Get()`时再次抛出。

[更多API文档](https://pkg.go.dev/github.com/timandy/routine#section-documentation)

# 垃圾回收

`routine`库内部维护了全局的`globalMap`变量，它存储了全部协程的上下文变量信息，在读写时基于协程的`goid`和协程变量的`ptr`进行变量寻址映射。

在进程的整个生命周期中，它可能会创建于销毁无数个协程，那么这些协程的上下文变量如何清理呢？

为解决这个问题，`routine`内部分配了一个全局的`gcTimer`，此定时器会在`globalMap`需要被清理时启动，定时扫描并清理`dead`协程在`globalMap`中缓存的上下文变量，从而避免可能出现的内存泄露隐患。

# License

MIT

# 鸣谢

这个库是从 [go-eden/routine](https://github.com/go-eden/routine) 分支出来的. 感谢 go-eden 的伟大工作!

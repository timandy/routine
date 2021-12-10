# routine

[![Build Status](https://github.com/timandy/routine/actions/workflows/build.yml/badge.svg)](https://github.com/timandy/routine/actions)
[![Codecov](https://codecov.io/gh/timandy/routine/branch/main/graph/badge.svg)](https://codecov.io/gh/timandy/routine)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/timandy/routine)

> [English Version](README_zh.md)

`routine`封装并提供了一些易用、高性能的`goroutine`上下文访问接口，它可以帮助你更优雅地访问协程上下文信息，但你也可能就此打开了潘多拉魔盒。

# 介绍

`Golang`语言从设计之初，就一直在不遗余力地向开发者屏蔽协程上下文的概念，包括协程`goid`的获取、进程内部协程状态、协程上下文存储等。

如果你使用过其他语言如`C++/Java`等，那么你一定很熟悉`ThreadLocal`，而在开始使用`Golang`之后，你一定会为缺少类似`ThreadLocal`的便捷功能而深感困惑与苦恼。 当然你可以选择使用`Context`
，让它携带着全部上下文信息，在所有函数的第一个输入参数中出现，然后在你的系统中到处穿梭。

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
	fmt.Printf("curr goid: %d\n", goid)
	fmt.Printf("all goids: %v\n", goids)
}
```

此例中`main`函数启动了一个新的协程，因此`Goid()`返回了主协程`1`，`AllGoids()`返回了主协程及协程`18`:

```text
curr goid: 1
all goids: [1 18]
```

## 使用`LocalStorage`

以下代码简单演示了`LocalStorage`的创建、设置、获取、跨协程传播等：

```go
package main

import (
	"fmt"
	"github.com/timandy/routine"
	"time"
)

var nameVar = routine.NewLocalStorage()

func main() {
	nameVar.Set("hello world")
	fmt.Println("name: ", nameVar.Get())

	// 其他协程不能读取前面Set的"hello world"
	go func() {
		fmt.Println("name1: ", nameVar.Get())
	}()

	// 但是可以通过Go函数启动新协程，并将当前main协程的全部协程上下文变量赋值过去
	routine.Go(func() {
		fmt.Println("name2: ", nameVar.Get())
	})

	// 或者，你也可以手动copy当前协程上下文至新协程，Go()函数的内部实现也是如此
	ic := routine.BackupContext()
	go func() {
		routine.RestoreContext(ic)
		fmt.Println("name3: ", nameVar.Get())
	}()

	time.Sleep(time.Second)
}
```

执行结果为：

```text
name:  hello world
name1:  <nil>
name3:  hello world
name2:  hello world
```

# API文档

此章节详细介绍了`routine`库封装的全部接口，以及它们的核心功能、实现方式等。

## `Goid() (id int64)`

获取当前`goroutine`的`goid`。

在正常情况下，`Goid()`优先尝试通过`go_tls`的方式直接获取，此操作性能极高，耗时通常只相当于`rand.Int()`的五分之一。

若出现版本不兼容等错误时，`Goid()`会尝试降级，即从`runtime.Stack`信息中解析获取，此时性能会急剧下降约千倍，但它可以保证功能正常可用。

## `AllGoids() (ids []int64)`

获取当前进程全部活跃`goroutine`的`goid`。

在`go 1.15`及更旧的版本中，`AllGoids()`会尝试从`runtime.Stack`信息中解析获取全部协程信息，但此操作非常低效，非常不建议在高频逻辑中使用。

在`go 1.16`之后的版本中，`AllGoids()`会通过`native`的方式直接读取`runtime`的全局协程池信息，在性能上得到了极大的提高， 但考虑到生产环境中可能有万、百万级的协程数量，因此仍不建议在高频使用它。

## `NewLocalStorage()`:

创建一个新的`LocalStorage`实例，它的设计思路与用法和其他语言中的`ThreadLocal`非常相似。

## `BackupContext() *ImmutableContext`

备份当前协程上下文的`local storage`数据，它只是一个便于上下文数据传递的不可变结构体。

## `RestoreContext(ic *ImmutableContext)`

主动继承备份到的上下文`local storage`数据，它会将其他协程`BackupContext()`的数据复制入当前协程上下文中，从而支持**跨协程的上下文数据传播**。

## `Go(f func())`

启动一个新的协程，同时自动将当前协程的全部上下文`local storage`数据复制至新协程，它的内部实现由`BackupContext()`和`RestoreContext()`组成。

## `LocalStorage`

表示协程上下文变量，支持的函数包括：

+ `Get() (value interface{})`：获取当前协程已设置的变量值，若未设置则为`nil`
+ `Set(v interface{}) interface{}`：设置当前协程的上下文变量值，返回之前已设置的旧值
+ `Del() (v interface{})`：删除当前协程的上下文变量值，返回已删除的旧值
+ `Clear()`：彻底清理此上下文变量在所有协程中保存的旧值

**提示：`Get/Set/Del`的内部实现采用无锁设计，在大部分情况下，它的性能表现都应该非常稳定且高效。**

# 垃圾回收

`routine`库内部维护了全局的`storages`变量，它存储了全部协程的上下文变量信息，在读写时基于协程的`goid`和协程变量的`ptr`进行变量寻址映射。

在进程的整个生命周期中，它可能会创建于销毁无数个协程，那么这些协程的上下文变量如何清理呢？

为解决这个问题，`routine`内部分配了一个全局的`GCTimer`，此定时器会在`storages`需要被清理时启动，定时扫描并清理`dead`协程在`storages`中缓存的上下文变量，从而避免可能出现的内存泄露隐患。

# License

MIT

# 鸣谢

这个库是从 [go-eden/routine](https://github.com/go-eden/routine) 分支出来的. 感谢 go-eden 的伟大工作!

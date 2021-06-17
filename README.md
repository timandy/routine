# routine

`routine` encapsulates and provides some easy-to-use, high-performance `goroutine` context access interfaces, which can
help you access coroutine context information more elegantly, but you may also open Pandora's Box.

# Introduce

The `Golang` language has been sparing no effort to shield developers from the concept of coroutine context from the
beginning of its design, including the acquisition of coroutine `goid`, the state of the coroutine within the process,
and the storage of coroutine context.

If you have used other languages such as `C++/Java/...`, then you must be familiar with `ThreadLocal`, and after
starting to use `Golang`, you will definitely feel confused and distressed by the lack of convenient functions similar
to `ThreadLocal` . Of course, you can choose to use `Context`, let it carry all the context information, appear in the
first input parameter of all functions, and then shuttle around in your system.

The core goal of `routine` is to open up another path: to introduce `goroutine local storage` into the world of `Golang`
, and at the same time expose the coroutine information to meet the needs of some people.

# Usage & Demo

This chapter briefly introduces how to install and use the `routine` library.

## Install

```bash
go get github.com/go-eden/routine
```

## Use `goid`

The following code simply demonstrates the use of `routine.Goid()` and `routine.AllGoids()`:

```go
package main

import (
	"fmt"
	"github.com/go-eden/routine"
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

In this example, the `main` function starts a new coroutine, so `Goid()` returns the main coroutine `1`,
and `AllGoids()` returns the main coroutine and coroutine `18`:

```text
curr goid: 1
all goids: [1 18]
```

## Use `LocalStorage`

The following code simply demonstrates the creation, setting, acquisition, and cross-coroutine propagation
of `LocalStorage`:

```go
package main

import (
	"fmt"
	"github.com/go-eden/routine"
	"time"
)

var nameVar = routine.NewLocalStorage()

func main() {
	nameVar.Set("hello world")
	fmt.Println("name: ", nameVar.Get())

	// other goroutine cannot read nameVar
	go func() {
		fmt.Println("name1: ", nameVar.Get())
	}()

	// but, the new goroutine could inherit/copy all local data from the current goroutine like this:
	routine.Go(func() {
		fmt.Println("name2: ", nameVar.Get())
	})

	// or, you could copy all local data manually
	ic := routine.BackupContext()
	go func() {
		routine.InheritContext(ic)
		fmt.Println("name3: ", nameVar.Get())
	}()

	time.Sleep(time.Second)
}
```

The results of the upper example are:

```text
name:  hello world
name1:  <nil>
name3:  hello world
name2:  hello world
```

# API

This chapter introduces in detail all the interfaces encapsulated by the `routine` library, as well as their core
functions and implementation methods.

## `Goid() (id int64)`

Get the `goid` of the current `goroutine`.

Under normal circumstances, `Goid()` first tries to get it directly through `go_tls`. This operation is extremely fast,
and the time-consuming is usually only one-fifth of `rand.Int()`.

If an error such as version incompatibility occurs, `Goid()` will try to parse it from the `runtime.Stack` information.
At this time, the performance will suffer exponential loss, which is about a thousand times slower, but the function can
be guaranteed to be normal Available.

## `AllGoids() (ids []int64)`

Get the `goid` of all active `goroutine` of the current process.

In `go 1.15` and older versions, `AllGoids()` will try to parse and get all the coroutine information from
the `runtime.Stack` information, but this operation is very inefficient and it is not recommended to use it in
high-frequency logic. .

In versions after `go 1.16`, `AllGoids()` will directly read the global coroutine pool information of `runtime`
through `native`, which has greatly improved performance, but considering the production environment There may be tens
of thousands or millions of coroutines, so it is still not recommended to use it at high frequencies.

## `NewLocalStorage()`:

Create a new instance of `LocalStorage`, its design idea is very similar to the usage of `ThreadLocal` in other
languages.

## `BackupContext() *ImmutableContext`

Back up the `local storage` data of the current coroutine context. It is just an immutable structure that facilitates
the transfer of context data.

## `InheritContext(ic *ImmutableContext)`

Actively inherit the backed-up context `local storage` data, it will copy the data of other coroutines `BackupContext()`
into the current coroutine context, thus supporting the contextual data propagation across coroutines.

## `Go(f func())`

Start a new coroutine and automatically copy all the context `local storage` data of the current coroutine to the new
coroutine. Its internal implementation consists of `BackupContext()` and `InheritContext()`.

## `LocalStorage`

Represents the context variable of the coroutine, and the supported functions include:

+ `Get() (value interface{})`: Get the variable value that has been set by the current coroutine.
+ `Set(v interface{}) interface{}`: Set the value of the context variable of the current coroutine, and return the old
  value that has been set before.
+ `Del() (v interface{})`: Delete the context variable value of the current coroutine and return the deleted old value.
+ `Clear()`: Thoroughly clean up the old value of this context variable saved in all coroutines.

# Garbage Collection

The `routine` library internally maintains the global `storages`, which stores all the variable values of all
coroutines, and performs data unique mapping based on the `goid` and `LocalStorage` of `goroutine` when reading and
writing.

In the entire life cycle of a process, there may be countless creation and destruction of coroutines, so it is necessary
to actively clean up the context data cached by the `dead` coroutine in the global `storages`. This work is performed by
a global timer in the `routine` library, which will, when necessary, Scan and clean up the relevant information of
the `dead` coroutine at regular intervals to avoid potential memory leaks.

# License

MIT
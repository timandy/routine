<!--变更日志-->

# v1.1.1 Release notes

### Features

- Support go version range `go1.13` ~ `go1.20`(New support `go1.20`).

### Changes

- Fix a memory leak risk caused by Timer.
- Reduce memory by using less chan.
- Update copyright.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.1.1](https://github.com/timandy/routine/tree/v1.1.1)

---

# v1.1.0 Release notes

### Features

- Support more arch `loong64`, `mips`, `mipsle`, `mips64`, `mips64le`, `ppc64le`, `riscv64`, `wasm`.

### Changes

- Upgrade dependencies to the latest version.
- Modify continuous integration script to support go1.19.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.1.0](https://github.com/timandy/routine/tree/v1.1.0)

---

# v1.0.9 Release notes

### Features

- Support arch `386` & `amd64` on `freebsd` and arch `ppc64` & `s390x` on `linux`.
- Support `Cancel()` and `GetWithTimeout()` methods for type `Future`.
- Support checking whether the tasks created by `GoWait(CancelRunnable)` and `GoWaitResult(CancelCallable)` methods are canceled.

### Changes

- Fix spell error of type `Future`.
- Rename type `Any` to `any`.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.9](https://github.com/timandy/routine/tree/v1.0.9)

---

# v1.0.8 Release notes

### Changes

- Rename `StackError` to `RuntimeError`.
- Support error nesting for `RuntimeError`.
- Beautify the error message of `RuntimeError`.
- Remove `bytesconv.Bytes()` and `bytesconv.String()` methods.
- Restore to the previous value if an overflow occurs when getting the index of `ThreadLocal`.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.8](https://github.com/timandy/routine/tree/v1.0.8)

---

# v1.0.7 Release notes

### Bugs

- Fix released `thread` struct may be resurrected from invalid memory and cause fault error.

### Changes

- Modify the error message format of `StackError`.
- Define function type `Supplier` for `threadLocal` and `inheritableThreadLocal` types.
- Define function type `Runnable` and `Callable` for `Go(Runnable)`, `GoWait(Runnable)` and `GoWaitResult(Callable)` methods.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.7](https://github.com/timandy/routine/tree/v1.0.7)

---

# v1.0.6 Release notes

### Bugs

- Fix fault error when pprof is running.

### Features

- Support more architectures `386`, `amd64`, `armv6`, `armv7`, `arm64`.

### Changes

- Read and write `coroutine` information through the `gohack` library, theoretically support unreleased `go` versions in the future.
- When `runtime.g` cannot be obtained natively, `panic` directly instead of falling back to invoke `runtime.Stack()` method.
- Remove api `ThreadLocal.Id()`.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.6](https://github.com/timandy/routine/tree/v1.0.6)

---

# v1.0.5 Release notes

### Features

- Support go version range `go1.13` ~ `go1.18`(New support `go1.18`).

### Changes

- Change license to `Apache-2.0`.
- Upgrade dependencies to the latest version.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.5](https://github.com/timandy/routine/tree/v1.0.5)

---

# v1.0.4 Release notes

### Features

- Add zero-copy conversion method between `bytes` and `string`, see `bytesconv.Bytes()` and `bytesconv.String()` methods.

### Changes

- Modify the garbage collection mechanism, remove `gcTimer`, no longer perform garbage collection through timers.
- Store the context in the `g.labels` field of the coroutine structure which will be set to `nil` after coroutine ends. The context data will be collected at the next `GC`.
- Use `go:linkname` to invoke assembly code `getg()` directly to improve performance.
- Implement the `getGoidByStack()` method by invoke `http.http2curGoroutineID()`.
- Remove api `AllGoids()` and `ForeachGoid()`.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.4](https://github.com/timandy/routine/tree/v1.0.4)

---

# v1.0.3 Release notes

### Features

- Support copy `Cloneable` objects to sub goroutine when create sub goroutines by `Go()`, `GoWait()` and `GoWaitResult()` methods.
- Add api `ForeachGoid(func(goid int64))` to run a func for each goid.

### Changes

- Support go version range `go1.13` ~ `go1.17`(Not support `go1.12` anymore).
- Use segment locks to reduce competition and improve `ThreadLocal`'s `read`, `write` and `gc` performance.
- Get all goids through `runtime.allgs` instead of `runtime.atomicAllG`, so `go1.13` ~ `go1.15` can also get all goids natively.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.3](https://github.com/timandy/routine/tree/v1.0.3)

---

# v1.0.2 Release notes

### Bugs

- Fix bug in `getAllGoidByStack()` method, Buffer may too small when dump all stack info.

### Features

- Support initialize value when first get from `ThreadLocal`.
- Add `StackError` to catch stack info.
- Add `Feature` to wait goroutine finished or get result from goroutine.
- Add api `NewThreadLocalWithInitial()`, `NewInheritableThreadLocal()` and `NewInheritableThreadLocalWithInitial()`.
- Support Inherit values of `ThreadLocal` by `Go`, `GoWait()` and `GoWaitResult()`.

### Changes

- Rename `LocalStorage` to `ThreadLocal`.
- Remove api `Clear()`, `InheritContext()` and `RestoreContext()`.
- Improve `gc` performance by reducing the number of for loops.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.2](https://github.com/timandy/routine/tree/v1.0.2)

---

# v1.0.1 Release notes

### Features

- Improve performance by use slice to store goroutine local values.
- Optimize `clearDeadStore()` method.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.1](https://github.com/timandy/routine/tree/v1.0.1)

---

# v1.0.0 Release notes

`This is the first stable version available for production. It is highly recommended to upgrade to this version if you have used a previous version.`

### Bugs

- Fix `NewLocalStorage()` always return the same value, so we can define multi `LocalStorage` instances.
- Fix `NewLocalStorage()` clear other `LocalStorage`'s value.
- Fix `RestoreContext()` not clear values when restore from empty `*ImmutableContext`.

### Features

- Not force create `store` when invoke `Get()`, `Remove()`, `Clear()`, `BackupContext()` methods to reduce memory usage.

### Changes

- Rename `InheritContext()` to `RestoreContext()`.
- Rename `Del()` to `Remove()`.
- Move Clear() method to `routine` package.

# Links

- Source code [https://github.com/timandy/routine/tree/v1.0.0](https://github.com/timandy/routine/tree/v1.0.0)

---

# v0.0.2 Release notes

### Features

- Support go version range `go1.12` ~ `go1.17`(New support `go1.17`).
- Enable GitHub actions for continuous integration.

### Known Issues

- `NewLocalStorage()` always return the same value.

# Links

- Source code [https://github.com/timandy/routine/tree/v0.0.2](https://github.com/timandy/routine/tree/v0.0.2)

---

# v0.0.1 Release notes

### Features

- Support go version range `go1.12` ~ `go1.16`.
- Support `Goid()` to get current goroutine id.
- Support `AllGoids` to get all goroutine ids.
- Support `ThreadLocal` to save values ingo to goroutine.

### Known Issues

- `NewLocalStorage()` always return the same value.

# Links

- Source code [https://github.com/timandy/routine/tree/v0.0.1](https://github.com/timandy/routine/tree/v0.0.1)

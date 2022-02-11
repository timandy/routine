<!--变更日志-->

# v1.0.3 Release notes

### Features

- Support copy `Cloneable` objects to sub goroutine when create sub goroutines by `Go()`,`GoWait()` and `GoWaitResult()` methods.
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

- Not force create `store` when invoke `Get()`,`Remove()`,`Clear()`,`BackupContext()` methods to reduce memory usage.

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
- Enable github actions for continuous integration.

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

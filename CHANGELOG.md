<!--变更日志-->

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

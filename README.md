# go_exec
A Go package providing an advanced API to execute commands

## Error Codes Contract

`pkg/remote` returns wrapped errors using `github.com/gvallee/go_errs/pkg/goerrs`.

Stable machine-readable codes currently used:

- `invalid_input`: caller-provided arguments are invalid (for example empty host, host with whitespace, empty binary path).
- `unavailable`: required dependency is unavailable (for example `ssh` binary not found).
- `internal`: command execution failed after validation passed.

### Usage

```go
res := remote.ExecCmd(host, binPath, args, env)
if res.Err != nil {
	switch {
	case goerrs.IsCode(res.Err, "invalid_input"):
		// treat as caller error
	case goerrs.IsCode(res.Err, "unavailable"):
		// retry or surface dependency issue
	case goerrs.IsCode(res.Err, "internal"):
		// unexpected runtime failure
	}
}
```

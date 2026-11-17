# Contributing

Thanks for helping! Issues and PRs are welcome in English or Spanish.

## Branching

- `main` is the production branch. It is protected: no direct pushes and no force-pushes, changes arrive through pull requests.
- `dev` is the integration branch. Branch from `dev` (`feat/...`, `fix/...`), open a PR to `dev`, and `dev` is merged into `main` through a PR for releases.
- Merged branches are deleted automatically.
- Copilot code review runs on PRs targeting `main`.

## Development

Requirements: Windows, Go 1.26+.

```bash
go vet ./...
go test ./...
gofmt -l .          # must print nothing
go build -o dist/taskglass.exe ./cmd/taskglass
powershell -File scripts/build-tap.ps1     # needs MinGW-w64 (winget install MartinStorsjo.LLVM-MinGW.UCRT)
go run ./cmd/taskglass --verbose
```

Rules of thumb:

- Keep Win32 calls inside `internal/winapi`. Put decisions in `internal/state` and `internal/config` so they stay unit-testable.
- Add or update tests for behaviour changes; update `docs/` and `CHANGELOG.md`.
- The native TAP runs inside `explorer.exe`: keep it small, check every HRESULT, never touch XAML objects off the UI thread, and never crash. Test changes on a real machine and mention the Windows build.
- Never leave the taskbar modified on a failure path: anything that applies an effect must be undone by `Restore`.
- Conventional-style commit messages (`feat:`, `fix:`, `docs:`...) are appreciated.

## Knowledge graph (optional)

The repo can be explored with [graphify](https://github.com/safishamsi/graphify): `graphify update .` rebuilds `graphify-out/` from code with no LLM. The output is committed so contributors and agents can query it (`graphify query "..."`).

# k9sight

A fast, keyboard-driven TUI for debugging Kubernetes workloads.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## Features

- Browse deployments, statefulsets, daemonsets, jobs, cronjobs
- View pod logs with error highlighting
- Monitor events and resource metrics
- Debug helpers for common issues (CrashLoopBackOff, ImagePullBackOff, etc.)
- Vim-style navigation
- Live filtering

## Install

```bash
go install github.com/doganarif/k9sight/cmd/k9sight@latest
```

Or build from source:

```bash
git clone https://github.com/doganarif/k9sight.git
cd k9sight
go build -o k9sight ./cmd/k9sight
```

## Usage

```bash
k9sight
```

### Key Bindings

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Navigate |
| `enter` | Select |
| `esc` | Back |
| `/` | Filter |
| `n` | Change namespace |
| `t` | Change resource type |
| `1-4` | Focus panel |
| `v` | Fullscreen panel |
| `?` | Help |
| `q` | Quit |

## Requirements

- Go 1.21+
- kubectl configured with cluster access

## License

MIT - see [LICENSE](LICENSE)

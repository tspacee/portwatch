# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes with configurable rules.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
rules:
  - port: 22
    allowed: true
  - port: 8080
    allowed: true
  - range: "1024-65535"
    allowed: false
    alert: true
```

When an unexpected port opens or closes, `portwatch` logs an alert and optionally executes a custom webhook or shell command.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `30s` | Poll interval |
| `--verbose` | `false` | Enable verbose logging |

---

## License

[MIT](LICENSE)
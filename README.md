# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

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

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Take a one-time snapshot of currently open ports:

```bash
portwatch snapshot
```

When an unexpected port opens or closes, portwatch prints an alert to stdout (and optionally to a log file):

```
[ALERT] New port detected: TCP 8080 (pid: 3421, process: python3)
[ALERT] Port closed: TCP 3306 (previously open)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--log` | `` | Path to log file (optional) |
| `--allowlist` | `` | Path to allowlist config file |

Define an allowlist to suppress alerts for known ports:

```bash
portwatch start --allowlist ./allowlist.yaml
```

---

## License

MIT © 2024 yourusername
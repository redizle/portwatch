# portwatch

A lightweight CLI daemon that monitors and logs port activity on localhost with alerting hooks.

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

Start watching specific ports:

```bash
portwatch --ports 8080,5432,6379
```

Run as a daemon with a log file and alert webhook:

```bash
portwatch --ports 80,443 --log /var/log/portwatch.log --webhook https://hooks.example.com/alert
```

**Example output:**

```
[2024-01-15 10:32:01] OPEN   port 8080 — process: main (pid 4821)
[2024-01-15 10:33:44] CLOSED port 8080 — process: main (pid 4821)
[2024-01-15 10:35:12] OPEN   port 5432 — process: postgres (pid 312)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ports` | `80,443` | Comma-separated list of ports to watch |
| `--interval` | `5s` | Polling interval |
| `--log` | stdout | Path to log file |
| `--webhook` | — | URL to POST alerts to |

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)
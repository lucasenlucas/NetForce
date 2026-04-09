# NetForce

```
  _   _      _   _____
 | \ | |    | | |  ___|
 |  \| | ___| |_| |_ ___  _ __ ___ ___
 | . ^ |/ _ \ __|  _/ _ \| '__/ __/ _ \
 | |\  |  __/ |_| || (_) | | | (_|  __/
 |_| \_|\___|\__\_| \___/|_|  \___\___|
```

**Performance & Resilience Testing Tool for Websites and Web Services.**

Part of the NetSuite ecosystem — consistent tooling for developers, students and admins.

---

> ⚠️ **AUTHORIZED USE ONLY**
> NetForce is designed strictly for systems you own or have explicit written permission to test.
> Unauthorized use is illegal and unethical. This tool is not for offensive use.

---

## ⚡ Quick Install

**One command — requires [Go 1.21+](https://golang.org/dl/)**

```bash
go install github.com/lucasenlucas/netforce/cmd/netforce@latest
```

After install, run it anywhere:
```bash
netforce -f explain
```

---

## Manual Installation (build from source)

```bash
git clone https://github.com/lucasenlucas/NetForce.git
cd NetForce
go mod tidy
go build -o netforce ./cmd/netforce
./netforce -f explain
```

---

## CLI Style

NetForce follows the same flag philosophy as **NetScope**:

```
netforce -d <domain> -f <feature> [options]
```

| Flag | Description |
|---|---|
| `-d, --domain` | Target domain or URL |
| `-f, --feature` | Feature to run (see below) |
| `-r, --rate` | Requests per second (default: 10) |
| `-t, --threads` | Concurrent workers (default: 5) |
| `--duration` | Test duration in seconds (default: 10) |
| `--timeout` | Request timeout per connection (default: 10s) |
| `--path` | URL endpoint path (default: /) |
| `--https` | Force HTTPS |
| `--safe` | Enable safe mode caps |
| `--output` | Output format: simple \| detailed \| json |
| `--report` | Save report to file after test |
| `--detect-rate-limit` | Detect HTTP 429 rate limiting |
| `--analyze-performance` | Analyze response time degradation |

---

## Features

### `-f stress` — Constant Load Test
Sends a steady, continuous stream of requests.
Best for understanding baseline performance.

```bash
netforce -d example.com -f stress -r 50 -t 5 --duration 30
```

### `-f ramp` — Gradual Load Increase *(coming soon)*
Slowly increases traffic from low to high.
Useful for finding the server's breaking point.

```bash
netforce -d example.com -f ramp -r 100 --duration 60 --report
```

### `-f spike` — Sudden Burst *(coming soon)*
Sends a large burst of traffic for a short time.
Simulates a flash sale or viral post.

```bash
netforce -d example.com -f spike -r 500 --duration 20
```

### `-f pulse` — Repeating Waves *(coming soon)*
Alternates between high and low traffic.
Models real-world peak patterns.

```bash
netforce -d example.com -f pulse -r 80 --duration 60
```

### `-f quick` — Beginner Safe Test
Runs a very gentle test using built-in safe defaults.
Great for your first look at server behavior.

```bash
netforce -d example.com -f quick
```

### `-f explain` — Plain English Explanation
Explains what NetForce does in simple language.
No target required.

```bash
netforce -f explain
```

---

## Output Examples

### Simple (default)
```
╔══════════════════════════════════════════╗
║       NetForce — Test Results            ║
╚══════════════════════════════════════════╝

  Target:                example.com
  Feature:               stress  (mode: constant)
  Configured Rate:       50 req/s
  Test Duration:         30.0s

  Total Requests:        1482
  Successes:             1479
  Errors:                3
  Success Rate:          99.80%
  Error Rate:            0.20%

  Avg Latency:           34ms
  Max Latency:           290ms
```

### JSON output
```bash
netforce -d example.com -f stress -r 30 --duration 10 --output json
```

---

## Safe Testing

NetForce is built with safety-first defaults:

- **Quick mode** always uses minimal safe limits
- **`--safe`** caps rate, threads and duration for cautious testing
- **Confirmation prompt** always appears before any test starts
- All help text and output discourages unauthorized use

---

## Roadmap

- [ ] `-f ramp` — Full gradual load implementation
- [ ] `-f spike` — Full spike burst implementation
- [ ] `-f pulse` — Full pulse wave implementation
- [ ] `--live` — Live refreshing terminal stats
- [ ] `--output detailed` — Per-status breakdown + p50/p90/p99 latency
- [ ] Multi-path testing support
- [ ] POST method with payload support

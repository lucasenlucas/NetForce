<p align="center">
  <img src="https://github.com/lucasenlucas/lucas_cdn/blob/main/Scherm%C2%ADafbeelding%202026-04-09%20om%2017.34.29.png?raw=true" alt="NetForce Banner"/>
</p>

<p align="center">
  <strong>NetForce</strong> — Performance & Resilience Testing CLI
</p>

<p align="center">
  Stress • Load • Behavior — test your system before it breaks
</p>

<p align="center">
  <strong>Know your limits. Before your users do.</strong>
</p>

---

## ⚡ What is NetForce?

NetForce is a **high-performance load testing CLI** built to simulate real-world traffic on websites and APIs.

No dashboards. No setup hell.

Just:
`netforce` → push your system to its limits


Built as part of the **NetSuite ecosystem** — consistent tooling alongside NetScope.

---

## ⚠️ Authorized Use Only

NetForce is designed strictly for:

- Systems you own  
- Systems you have **explicit written permission** to test  

Unauthorized use is illegal. Don’t be that guy.

---

## 🚀 Quick Install

### macOS & Linux
```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/NetForce/main/install.sh | sh
```
Run instantly:
```bash
netforce -f explain
```
### Go install
```bash
go install github.com/lucasenlucas/netforce/cmd/netforce@latest
```

### Manual build
```bash
git clone https://github.com/lucasenlucas/NetForce.git
cd NetForce
go mod tidy
go build -o netforce ./cmd/netforce
./netforce -f explain
```

### CLI Philosophy
Same system as NetScope → no learning curve.
```bash
netforce -d <domain> -f <feature> [options]
```

### Core Flags

| Flag                    | Description              |
| ----------------------- | ------------------------ |
| `-d, --domain`          | Target domain or URL     |
| `-f, --feature`         | Feature to run           |
| `-r, --rate`            | Requests per second      |
| `-t, --threads`         | Concurrent workers       |
| `--duration`            | Test duration            |
| `--timeout`             | Request timeout          |
| `--path`                | Endpoint path            |
| `--https`               | Force HTTPS              |
| `--safe`                | Safe caps enabled        |
| `--output`              | simple \| detailed \| json |
| `--report`              | Save report              |
| `--detect-rate-limit`   | Detect HTTP 429          |
| `--analyze-performance` | Detect slowdown          |
| `--live`                | Live terminal stats      |

## Core Features
### stress — Constant Load
Push steady traffic → measure baseline performance
```bash
netforce -d example.com -f stress -r 50 -t 5 --duration 30
```
### ramp
Gradually increase load → find breaking point
```bash
netforce -d example.com -f ramp -r 100 --duration 60
```
### spike
Short aggressive burst → simulate viral traffic
```bash
netforce -d example.com -f spike -r 500 --duration 10
```
### pulse
Traffic waves → simulate real usage patterns
```bash
netforce -d example.com -f pulse -r 80 --duration 30
```
### quick — Safe Mode
Beginner-friendly, low-impact test
```bash
netforce -d example.com -f quick
```
### explain
Plain explanation of what NetForce does
```bash
netforce -f explain
```

## Output Default
```
╔══════════════════════════════════════════╗
║       NetForce — Test Results            ║
╚══════════════════════════════════════════╝

  Target:                example.com
  Feature:               stress
  Rate:                  50 req/s
  Duration:              30s

  Total Requests:        1482
  Success Rate:          99.80%
  Error Rate:            0.20%

  Avg Latency:           34ms
  Max Latency:           290ms
```

### Detailed & JSON Output
Get breakdown of status codes and latency percentiles (p50/p90/p99):
```bash
netforce -d example.com -f stress -r 30 --duration 10 --output detailed
netforce -d example.com -f stress -r 30 --duration 10 --output json
```

## Safety First
NetForce is built to protect users and systems:
* quick mode uses safe defaults
* --safe enforces strict caps
* Confirmation prompt before execution
* Clear warnings in CLI
(especially so that I also stay out of trouble with the law.)

## Why NetForce?
Most load tools are:
* Overcomplicated
* Slow to setup
* UI-heavy and bloated

NetForce is:
* Fast
* Focused
* Scriptable
* Built for real usage

## Roadmap
* Live terminal dashboard (--live) ✅
* Advanced latency stats (p50 / p90 / p99) ✅
* Multi-endpoint testing
* POST payload support

## Author
Built by Lucas Mangroelal 
https://lucasmangroelal.nl

❤️ Support
* Star the repo!
* Contribute!
* Share it!

## ⚠️ Disclaimer
This tool is for educational and authorized testing only.
Do not use against systems without permission.

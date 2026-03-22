# 🔋 cli-battery

[![GitHub release](https://img.shields.io/github/release/zoe668668/cli-battery.svg)](https://github.com/zoe668668/cli-battery/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/zoe668668/cli-battery.svg?style=social&label=Star)](https://github.com/zoe668668/cli-battery)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoe668668/cli-battery)](https://goreportcard.com/report/github.com/zoe668668/cli-battery)

> **A beautiful battery health monitor for your terminal** 🔋✨

A lightweight, zero-dependency CLI tool to monitor your laptop battery health, cycle count, temperature, and estimated lifespan with beautiful terminal graphics.

## ✨ Features

- 📊 **Battery Health Score** - Visual health percentage with color coding (green/yellow/red)
- 🔄 **Cycle Count** - Track charge cycles vs. design limit
- 🌡️ **Temperature Monitor** - Real-time battery temperature
- ⏱️ **Time Estimates** - Time to full charge / time remaining
- 📈 **Pretty Charts** - ASCII progress bars
- 🖥️ **Cross-Platform** - macOS & Linux support
- ⚡ **Zero Dependencies** - Single binary, no runtime required
- 🎨 **Color Themes** - Multiple color schemes (default, neon, dark, minimal)

## 📦 Installation

### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -sSL https://github.com/zoe668668/cli-battery/releases/latest/download/cli-battery-darwin-arm64 -o cli-battery
chmod +x cli-battery
sudo mv cli-battery /usr/local/bin/
```

### Linux

```bash
# x86_64
curl -sSL https://github.com/zoe668668/cli-battery/releases/latest/download/cli-battery-linux-amd64 -o cli-battery
chmod +x cli-battery
sudo mv cli-battery /usr/local/bin/
```

### Go Install

```bash
go install github.com/zoe668668/cli-battery@latest
```

## 🚀 Quick Start

```bash
# Basic usage - show battery status
cli-battery

# Watch mode - update every 5 seconds
cli-battery --watch

# JSON output for scripting
cli-battery --json

# Show with neon theme
cli-battery --theme neon
```

## 📸 Screenshots

### Basic Output

```
┌─────────────────────────────────────────────────────────┐
│  🔋 Battery Status                                       │
├─────────────────────────────────────────────────────────┤
│  Charge      ████████████████████░░░░░░░░░░  78%        │
│  Health      ████████████████████████████░░  94%        │
│  Cycles      ████████░░░░░░░░░░░░░░░░░░░░░░  287 / 1000 │
│  Temp        32.5°C                                      │
├─────────────────────────────────────────────────────────┤
│  ⏱️  Time to Full: 1h 23m                                │
│  🔌 Status: Charging (USB-C)                             │
│  📅 Estimated Life: ~4.2 years remaining                 │
└─────────────────────────────────────────────────────────┘
```

### Watch Mode

```bash
$ cli-battery --watch

[14:32:15] 🔋 78% | ⚡ Charging | 32°C | ETA: 1h 23m
[14:32:20] 🔋 79% | ⚡ Charging | 32°C | ETA: 1h 19m
[14:32:25] 🔋 80% | ⚡ Charging | 33°C | ETA: 1h 15m
```

### JSON Output

```bash
$ cli-battery --json
{
  "charge_percent": 78,
  "health_percent": 94,
  "cycle_count": 287,
  "max_cycles": 1000,
  "temperature_celsius": 32.5,
  "status": "charging",
  "power_source": "USB-C",
  "estimated_life_years": 4.2
}
```

## 🎨 Themes

| Theme | Description |
|-------|-------------|
| `default` | Standard terminal colors |
| `neon` | Vibrant neon colors |
| `dark` | No colors, text only |
| `minimal` | Clean minimal output |

```bash
cli-battery --theme neon
```

## 🔧 Options

```
Usage: cli-battery [options]

Options:
  -w, --watch         Watch mode, update every 5 seconds
  -i, --interval      Watch interval in seconds (default: 5)
  -j, --json          Output as JSON
  -d, --detail        Show detailed battery information
  -t, --theme         Color theme (default, dark, neon, minimal)
  -n, --no-color      Disable colored output
  -v, --version       Show version
  -h, --help          Show help
```

## 🛠️ Building from Source

```bash
git clone https://github.com/zoe668668/cli-battery.git
cd cli-battery
go build -o cli-battery .
./cli-battery
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [neofetch](https://github.com/dylanaraps/neofetch) and [btop](https://github.com/aristocratos/btop)
- Built with ❤️ using Go

---

⭐ **If this project helped you, please consider giving it a star!** ⭐

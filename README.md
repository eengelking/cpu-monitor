# CPU Monitor

A terminal-based CPU monitoring tool written in Go. Features real-time CPU usage visualization with per-core monitoring, memory tracking, and system metrics.

![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)
![License](https://img.shields.io/badge/License-MIT-green)

## Features

- 🚀 **Real-time CPU monitoring** with sub-second refresh rates
- 📊 **Per-core usage tracking** with multi-column layout for many-core systems
- 📈 **60-second history graph** with color-coded usage levels
- 💾 **Memory usage monitoring** with visual progress bars
- ⚡ **Low overhead** - optimized to use only 1-2% CPU
- 🎯 **Interactive controls** - pause, reset, and help system
- 🖥️ **Cross-platform** - works on macOS and Linux

## Installation

### Build from Source

#### Prerequisites

- Go 1.21 or higher
- Git

#### Quick Install

```bash
# Clone the repository
git clone https://github.com/yourusername/cpu-monitor.git
cd cpu-monitor

# Build and install
go build -o cpu-monitor ./cmd/cpu-monitor

# Optional: Move to PATH
sudo mv cpu-monitor /usr/local/bin/
```

## Building for Different Platforms

### macOS

#### Intel (amd64)
```bash
GOOS=darwin GOARCH=amd64 go build -o cpu-monitor-darwin-amd64 ./cmd/cpu-monitor
```

#### Apple Silicon (arm64)
```bash
GOOS=darwin GOARCH=arm64 go build -o cpu-monitor-darwin-arm64 ./cmd/cpu-monitor
```

### Linux

#### Intel/AMD (amd64)
```bash
GOOS=linux GOARCH=amd64 go build -o cpu-monitor-linux-amd64 ./cmd/cpu-monitor
```

#### ARM64
```bash
GOOS=linux GOARCH=arm64 go build -o cpu-monitor-linux-arm64 ./cmd/cpu-monitor
```

## Usage

### Basic Usage

```bash
# Run with default settings
cpu-monitor

# Show help
cpu-monitor -help
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-refresh ms` | Set refresh rate in milliseconds (100-5000) | 500 |
| `-history n` | Number of history points to keep | 120 |
| `-avg n` | Moving average window size | 10 |
| `-help` | Show command line help | - |

### Examples

```bash
# Update every second
cpu-monitor -refresh 1000

# Reduce history buffer for lower memory usage
cpu-monitor -history 60 -avg 5

# Fast refresh for detailed monitoring
cpu-monitor -refresh 100
```

## Keyboard Controls

| Key | Action |
|-----|--------|
| `h` | Toggle help screen |
| `q`, `Ctrl+C` | Quit application |
| `r` | Reset CPU history |
| `p` | Pause/unpause monitoring |

## Display Sections

### Main View
- **Total CPU**: Overall system CPU usage percentage
- **Moving Average**: 10-sample moving average of CPU usage
- **Core Bars**: Individual CPU core usage with htop-style bars
- **CPU History**: 60-second vertical bar graph showing usage over time
- **Memory**: System RAM usage with visual progress bar
- **System Info**: Load average, process count, and uptime

### Color Indicators
- 🟢 **Green** (0-30%): Low usage
- 🔵 **Blue** (30-50%): Light usage
- 🟡 **Yellow** (50-70%): Moderate usage
- 🟠 **Orange** (70-90%): High usage
- 🔴 **Red** (90-100%): Critical usage

## Performance

CPU Monitor is highly optimized for minimal system impact:
- **CPU Usage**: ~1-2% on modern systems
- **Memory**: ~15-25 MB RAM
- **Update Rate**: Configurable from 100ms to 5 seconds

### Optimization Features
- Cached static system information
- Throttled expensive operations (temperature, process count)
- Pre-rendered UI styles
- Efficient string building
- Proper CPU sampling intervals

## System Requirements

### Minimum Requirements
- **OS**: macOS 10.15+ or Linux kernel 3.0+
- **Terminal**: Any terminal with 256-color support
- **Width**: Minimum 80 columns
- **Height**: Minimum 24 rows

### Recommended
- **Terminal**: iTerm2, Alacritty, or Kitty
- **Font**: Monospace font with Unicode support
- **Size**: 120+ columns for optimal display

## Technical Details

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [gopsutil](https://github.com/shirou/gopsutil) - System information gathering

### Project Structure
```
cpu-monitor/
├── cmd/
│   └── cpu-monitor/    # Application entry point
├── internal/
│   ├── config/         # Configuration and constants
│   ├── metrics/        # System metrics collection
│   └── ui/             # Terminal UI components
├── go.mod              # Go module definition
└── go.sum              # Dependency checksums
```

## Troubleshooting

### High CPU Usage
If CPU Monitor itself is using too much CPU:
1. Increase the refresh rate: `cpu-monitor -refresh 1000`
2. Reduce history size: `cpu-monitor -history 60`
3. Check for other system monitoring tools running simultaneously

### Temperature Not Showing
- **macOS**: Temperature sensors may require additional permissions or tools
- **Linux**: Ensure `lm-sensors` is installed: `sudo apt-get install lm-sensors`

### Display Issues
- Ensure your terminal supports 256 colors: `echo $TERM`
- Try resizing your terminal window
- Use a monospace font with Unicode support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- System metrics via [gopsutil](https://github.com/shirou/gopsutil)
- Inspired by htop, btop, and other terminal monitoring tools

## Author

Created with ❤️ for the terminal by [Ed Engelking](https://github.com/eengelking) as an experiment with [Claude Code](https://www.anthropic.com/claude-code).
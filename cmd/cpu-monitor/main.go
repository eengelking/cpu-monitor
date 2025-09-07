package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/cpu-monitor/internal/config"
	"github.com/user/cpu-monitor/internal/ui"
)

func main() {
	var (
		refreshRate = flag.Int("refresh", 500, "Refresh rate in milliseconds (default: 500)")
		historySize = flag.Int("history", 120, "Number of history points to keep (default: 120)")
		avgSize     = flag.Int("avg", 10, "Moving average window size (default: 10)")
		help        = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	if *refreshRate < 100 {
		fmt.Println("Warning: Refresh rate too low, setting to minimum 100ms")
		*refreshRate = 100
	}

	if *refreshRate > 5000 {
		fmt.Println("Warning: Refresh rate too high, setting to maximum 5000ms")
		*refreshRate = 5000
	}

	cfg := config.Config{
		RefreshRate:   time.Duration(*refreshRate) * time.Millisecond,
		HistorySize:   *historySize,
		MovingAvgSize: *avgSize,
	}

	model := ui.NewModel(cfg)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	banner := `
╔═══════════════════════════════════════════════════════════════════════════════╗
║                           CPU MONITOR - HELP                                  ║
╚═══════════════════════════════════════════════════════════════════════════════╝

A cyberpunk-themed terminal CPU monitoring tool with real-time visualization.

USAGE:
    cpu-monitor [OPTIONS]

OPTIONS:
    -refresh <ms>    Set refresh rate in milliseconds (100-5000, default: 500)
    -history <n>     Number of history points to keep (default: 120)
    -avg <n>         Moving average window size (default: 10)
    -help            Show this help message

KEYBOARD CONTROLS:
    q, Ctrl+C        Quit the application
    r                Reset history
    p                Pause/unpause monitoring

FEATURES:
    • Real-time CPU usage monitoring with per-core breakdown
    • Visual ASCII graphs showing CPU history
    • Memory usage tracking
    • System load average display
    • Process count monitoring
    • CPU temperature display (when available)
    • Animated cyberpunk aesthetic with neon colors

VISUAL INDICATORS:
    Green  (< 30%)   Low CPU usage
    Yellow (30-70%)  Moderate CPU usage
    Red    (> 70%)   High CPU usage

EXAMPLES:
    cpu-monitor                      # Run with default settings
    cpu-monitor -refresh 1000        # Update every second
    cpu-monitor -history 60 -avg 5   # Keep 60 history points, 5-point average

Created with ♥ for the terminal
`
	fmt.Println(banner)
}
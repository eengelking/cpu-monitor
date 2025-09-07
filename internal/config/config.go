package config

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	RefreshRate   time.Duration
	HistorySize   int
	MovingAvgSize int
}

var DefaultConfig = Config{
	RefreshRate:   500 * time.Millisecond,
	HistorySize:   120,
	MovingAvgSize: 10,
}

type ColorScheme struct {
	Background   lipgloss.Color
	Foreground   lipgloss.Color
	NeonGreen    lipgloss.Color
	NeonBlue     lipgloss.Color
	NeonPurple   lipgloss.Color
	NeonPink     lipgloss.Color
	Yellow       lipgloss.Color
	Red          lipgloss.Color
	DimGray      lipgloss.Color
	BrightWhite  lipgloss.Color
}

var Colors = ColorScheme{
	Background:   lipgloss.Color("#000000"),
	Foreground:   lipgloss.Color("#FFFFFF"),
	NeonGreen:    lipgloss.Color("#00FF00"),
	NeonBlue:     lipgloss.Color("#00FFFF"),
	NeonPurple:   lipgloss.Color("#FF00FF"),
	NeonPink:     lipgloss.Color("#FF1493"),
	Yellow:       lipgloss.Color("#FFFF00"),
	Red:          lipgloss.Color("#FF0000"),
	DimGray:      lipgloss.Color("#404040"),
	BrightWhite:  lipgloss.Color("#FFFFFF"),
}

const (
	LowCPUThreshold  = 30.0
	HighCPUThreshold = 70.0
)

func GetCPUColor(usage float64) lipgloss.Color {
	switch {
	case usage < LowCPUThreshold:
		return Colors.NeonGreen
	case usage < HighCPUThreshold:
		return Colors.Yellow
	default:
		return Colors.Red
	}
}

const (
	BarFull   = "█"
	BarThree  = "▓"
	BarTwo    = "▒"
	BarOne    = "░"
	BarEmpty  = " "
)

const (
	GraphBar1 = "▁"
	GraphBar2 = "▂"
	GraphBar3 = "▃"
	GraphBar4 = "▄"
	GraphBar5 = "▅"
	GraphBar6 = "▆"
	GraphBar7 = "▇"
	GraphBar8 = "█"
)

const (
	BorderTopLeft     = "┌"
	BorderTopRight    = "┐"
	BorderBottomLeft  = "└"
	BorderBottomRight = "┘"
	BorderHorizontal  = "─"
	BorderVertical    = "│"
	BorderCross       = "┼"
	BorderTeeRight    = "├"
	BorderTeeLeft     = "┤"
	BorderTeeTop      = "┴"
	BorderTeeBottom   = "┬"
)

var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const AppTitle = "CPU Monitor"
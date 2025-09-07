package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/user/cpu-monitor/internal/config"
)

// Pre-created styles to avoid recreation on every render
var (
	TitleStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonGreen).
		Bold(true)

	HelpStyle = lipgloss.NewStyle().
		Foreground(config.Colors.DimGray)

	KeyStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue)

	BorderStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple)

	SpinnerStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple)

	PauseStyle = lipgloss.NewStyle().
		Foreground(config.Colors.Yellow).
		Bold(true).
		Blink(true)

	TimeStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPink)

	LoadStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonGreen)

	ProcessStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPink)

	UptimeStyle = lipgloss.NewStyle().
		Foreground(config.Colors.Yellow)

	DimGrayStyle = lipgloss.NewStyle().
		Foreground(config.Colors.DimGray)

	ScaleStyle = lipgloss.NewStyle().
		Foreground(config.Colors.DimGray)

	BracketStyle = lipgloss.NewStyle().
		Foreground(config.Colors.DimGray)

	GraphBorderStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue)

	GraphTitleStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple).
		Bold(true)
)

// CPU label styles by width
var (
	CPULabelStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue).
		Width(12).
		Align(lipgloss.Left)

	CPULabelStyleCompact = lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue).
		Width(5).
		Align(lipgloss.Left)

	MemoryLabelStyle = lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple).
		Width(12).
		Align(lipgloss.Left)
)

// Color styles for different CPU usage levels
var (
	GreenStyle  = lipgloss.NewStyle().Foreground(config.Colors.NeonGreen)
	BlueStyle   = lipgloss.NewStyle().Foreground(config.Colors.NeonBlue)
	YellowStyle = lipgloss.NewStyle().Foreground(config.Colors.Yellow)
	OrangeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8800"))
	RedStyle    = lipgloss.NewStyle().Foreground(config.Colors.Red)
)

func GetColorStyle(percentage float64) lipgloss.Style {
	switch {
	case percentage < 30:
		return GreenStyle
	case percentage < 50:
		return BlueStyle
	case percentage < 70:
		return YellowStyle
	case percentage < 90:
		return OrangeStyle
	default:
		return RedStyle
	}
}
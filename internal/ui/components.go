package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/cpu-monitor/internal/config"
)

// Pre-rendered brackets for progress bars
var (
	openBracket  = "["
	closeBracket = "]"
)

func CreateProgressBar(percentage float64, width int, color lipgloss.Color) string {
	if width <= 0 {
		return ""
	}
	
	// Adjust width for brackets
	barWidth := width - 2
	if barWidth <= 0 {
		return openBracket + closeBracket
	}

	fillWidth := int(float64(barWidth) * percentage / 100)
	if fillWidth > barWidth {
		fillWidth = barWidth
	}

	var bar strings.Builder
	
	// Build the bar with color transitions
	for i := 0; i < barWidth; i++ {
		if i < fillWidth {
			// Calculate percentage at this position
			positionPercentage := (float64(i) / float64(barWidth)) * 100
			
			// Use pre-cached color styles
			var style lipgloss.Style
			if positionPercentage < 30 {
				style = GreenStyle
			} else if positionPercentage < 50 {
				style = BlueStyle
			} else if positionPercentage < 70 {
				style = YellowStyle
			} else if positionPercentage < 90 {
				style = OrangeStyle
			} else {
				style = RedStyle
			}
			
			bar.WriteString(style.Render("|"))
		} else {
			// Empty space
			bar.WriteString(" ")
		}
	}

	return openBracket + bar.String() + closeBracket
}

func CreateASCIIGraph(values []float64, width, height int) string {
	if len(values) == 0 || width <= 0 || height <= 0 {
		return ""
	}

	// Reserve space for scale labels on the left (5 chars: "100%|")
	scaleWidth := 5
	graphWidth := width - scaleWidth
	if graphWidth < 20 {
		graphWidth = 20
	}

	// Initialize graph with empty spaces
	graph := make([][]string, height)
	for i := range graph {
		graph[i] = make([]string, graphWidth)
		for j := range graph[i] {
			graph[i][j] = " "
		}
	}

	// Determine how many data points to show
	startIdx := 0
	if len(values) > graphWidth {
		startIdx = len(values) - graphWidth
	}
	
	// Plot each value as a vertical bar (like rotated CPU bars)
	for x := 0; x < graphWidth && startIdx+x < len(values); x++ {
		value := values[startIdx+x]
		if value <= 0 {
			continue
		}
		
		// Calculate bar height
		barHeight := int(value * float64(height) / 100)
		if barHeight > height {
			barHeight = height
		}
		if barHeight == 0 && value > 0 {
			barHeight = 1
		}

		// Draw the vertical bar with color gradient (like horizontal bars)
		for y := 0; y < height; y++ {
			yPos := height - 1 - y
			
			if y < barHeight {
				// Calculate position percentage for color gradient
				positionPercentage := (float64(y) / float64(height)) * 100
				
				// Use cached style
				style := GetColorStyle(positionPercentage)
				graph[yPos][x] = style.Render("│")
			}
		}
	}

	// Build the result with scale on the left
	var result strings.Builder
	
	for i, row := range graph {
		// Add scale label every few rows
		scaleLabel := "    "
		if i == 0 {
			scaleLabel = "100%"
		} else if i == height/4 {
			scaleLabel = " 75%"
		} else if i == height/2 {
			scaleLabel = " 50%"
		} else if i == 3*height/4 {
			scaleLabel = " 25%"
		} else if i == height-1 {
			scaleLabel = "  0%"
		}
		
		result.WriteString(ScaleStyle.Render(scaleLabel))
		result.WriteString("│")
		
		for _, cell := range row {
			result.WriteString(cell)
		}
		result.WriteString("\n")
	}
	
	// Add bottom axis with brackets like the horizontal bars
	result.WriteString(ScaleStyle.Render("    "))
	result.WriteString(BracketStyle.Render("└" + strings.Repeat("─", graphWidth)))
	result.WriteString("\n")
	result.WriteString(ScaleStyle.Render("     "))
	result.WriteString(BracketStyle.Render("└60s" + strings.Repeat(" ", graphWidth/2-4) + "Now┘"))

	return result.String()
}

func CreateBox(title, content string, width, height int, borderColor lipgloss.Color) string {
	style := lipgloss.NewStyle().
		Foreground(borderColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Height(height).
		Padding(0, 1)

	if title != "" {
		titleStyle := lipgloss.NewStyle().
			Foreground(config.Colors.NeonBlue).
			Bold(true)
		
		titleLine := config.BorderTeeLeft + " " + titleStyle.Render(title) + " " + config.BorderTeeRight
		content = titleLine + "\n" + content
	}

	return style.Render(content)
}

func CreateSpinner(frame int) string {
	return SpinnerStyle.Render(config.SpinnerFrames[frame%len(config.SpinnerFrames)])
}

func CreateCPUBar(label string, percentage float64, width int) string {
	color := config.GetCPUColor(percentage)
	percentStr := fmt.Sprintf("%5.1f%%", percentage)
	
	// For compact display, use shorter labels
	compactLabel := label
	if width < 50 {
		// Shorten "Core X" to just the number for compact display
		if strings.HasPrefix(label, "Core ") {
			compactLabel = strings.TrimPrefix(label, "Core ")
			// Pad to 2 characters for alignment
			if len(compactLabel) == 1 {
				compactLabel = " " + compactLabel
			}
		} else if label == "Total CPU" {
			compactLabel = "Total"
		} else if label == "Moving Avg" {
			compactLabel = "Avg"
		}
	}
	
	labelWidth := 12
	if width < 50 {
		labelWidth = 5
	}
	
	labelStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue).
		Width(labelWidth).
		Align(lipgloss.Left)
	
	percentStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Width(7).
		Align(lipgloss.Right)
	
	// Calculate bar width (including brackets)
	barWidth := width - labelWidth - 7 - 2
	if barWidth < 12 {
		barWidth = 12
	}
	
	bar := CreateProgressBar(percentage, barWidth, color)
	
	return labelStyle.Render(compactLabel) + " " + bar + " " + percentStyle.Render(percentStr)
}

func CreateMemoryBar(used, total uint64, percentage float64, width int) string {
	color := config.GetCPUColor(percentage)
	usedStr := formatBytes(used)
	totalStr := formatBytes(total)
	percentStr := fmt.Sprintf("%6.2f%%", percentage)
	
	labelStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple).
		Width(12).
		Align(lipgloss.Left)
	
	infoStyle := lipgloss.NewStyle().
		Foreground(color).
		Width(20).
		Align(lipgloss.Right)
	
	percentStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Width(8).
		Align(lipgloss.Right)
	
	// Calculate bar width: total width - label(12) - info(20) - percentage(8) - spacing(3)
	barWidth := width - 43
	if barWidth < 22 {
		barWidth = 22
	}
	
	bar := CreateProgressBar(percentage, barWidth, color)
	info := fmt.Sprintf("%s/%s", usedStr, totalStr)
	
	return labelStyle.Render("Memory:") + " " + bar + " " + infoStyle.Render(info) + " " + percentStyle.Render(percentStr)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	exp := 0
	val := float64(bytes)
	
	for val >= unit && exp < len(units)-1 {
		val /= unit
		exp++
	}
	
	if exp == 0 {
		return fmt.Sprintf("%.0f %s", val, units[exp])
	}
	return fmt.Sprintf("%.1f %s", val, units[exp])
}


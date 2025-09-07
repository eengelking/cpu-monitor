package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/cpu-monitor/internal/config"
)

func (m Model) View() string {
	if m.metrics == nil {
		return "Initializing..."
	}

	// Show help screen if toggled
	if m.showHelp {
		return m.renderHelpScreen()
	}

	var b strings.Builder
	
	// Build the complete view
	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderSystemInfo())
	b.WriteString("\n\n")
	b.WriteString(m.renderCPUBars())
	b.WriteString("\n")
	b.WriteString(m.renderGraph())
	b.WriteString("\n")
	b.WriteString(m.renderMemoryInfo())
	b.WriteString("\n\n")
	b.WriteString(m.renderBottomInfo())
	
	return b.String()
}

func (m Model) renderHeader() string {
	title := TitleStyle.Render(config.AppTitle)
	
	controls := []string{
		KeyStyle.Render("h") + HelpStyle.Render(":help"),
		KeyStyle.Render("q") + HelpStyle.Render(":quit"),
		KeyStyle.Render("r") + HelpStyle.Render(":reset"),
		KeyStyle.Render("p") + HelpStyle.Render(":pause"),
	}
	
	controlsText := strings.Join(controls, "  ")
	
	// Calculate spacing to right-align controls
	titleWidth := lipgloss.Width(title)
	controlsWidth := lipgloss.Width(controlsText)
	spacing := m.width - titleWidth - controlsWidth
	if spacing < 2 {
		spacing = 2
	}
	
	header := title + strings.Repeat(" ", spacing) + controlsText
	
	separator := BorderStyle.Render(strings.Repeat("═", m.width))
	
	return header + "\n" + separator
}

func (m Model) renderSystemInfo() string {
	statusIndicator := CreateSpinner(m.spinnerFrame)
	if m.paused {
		statusIndicator = PauseStyle.Render("⏸ PAUSED")
	}
	
	currentTime := time.Now().Format("15:04:05.000")
	
	tempStr := "N/A"
	if m.metrics.Temperature > 0 {
		tempStyle := GetColorStyle(m.metrics.Temperature)
		tempStr = tempStyle.Render(fmt.Sprintf("%.1f°C", m.metrics.Temperature))
	}

	info := fmt.Sprintf(
		"%s  CPU: %s  Cores: %d  Threads: %d  Freq: %.0f MHz  Temp: %s  %s",
		statusIndicator,
		truncateString(m.metrics.ModelName, 20),
		m.metrics.CoreCount,
		m.metrics.ThreadCount,
		m.metrics.Frequency,
		tempStr,
		TimeStyle.Render(currentTime),
	)

	return info
}

func (m Model) renderCPUBars() string {
	barWidth := m.width

	var bars []string

	// Show total and moving average at the top
	totalBar := CreateCPUBar("Total CPU", m.metrics.TotalUsage, barWidth)
	bars = append(bars, totalBar)

	movingAvg := m.history.GetMovingAverage()
	avgBar := CreateCPUBar("Moving Avg", movingAvg, barWidth)
	bars = append(bars, avgBar)

	bars = append(bars, "")

	// Calculate layout for per-core display
	numCores := len(m.metrics.PerCoreUsage)
	if numCores == 0 {
		return strings.Join(bars, "\n")
	}

	// Determine number of columns based on terminal width
	// Each column needs roughly 40 characters minimum
	minColumnWidth := 40
	maxColumns := m.width / minColumnWidth
	if maxColumns < 1 {
		maxColumns = 1
	}
	if maxColumns > 4 {
		maxColumns = 4 // Cap at 4 columns for readability
	}

	// Calculate actual column width
	columnWidth := m.width / maxColumns
	
	// Calculate rows needed
	rowsNeeded := (numCores + maxColumns - 1) / maxColumns

	// Build the multi-column layout
	for row := 0; row < rowsNeeded; row++ {
		var rowBars []string
		
		for col := 0; col < maxColumns; col++ {
			coreIndex := row + col*rowsNeeded
			if coreIndex >= numCores {
				// Add empty space for alignment
				rowBars = append(rowBars, strings.Repeat(" ", columnWidth-1))
			} else {
				label := fmt.Sprintf("Core %d", coreIndex)
				bar := CreateCPUBar(label, m.metrics.PerCoreUsage[coreIndex], columnWidth-1)
				rowBars = append(rowBars, bar)
			}
		}
		
		bars = append(bars, strings.Join(rowBars, " "))
	}

	return strings.Join(bars, "\n")
}

func (m Model) renderGraph() string {
	graphWidth := m.width
	graphHeight := 8

	titleStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple).
		Bold(true)

	titleText := " CPU History (60s) "
	title := titleStyle.Render(titleText)
	titleWidth := lipgloss.Width(titleText)

	historyValues := m.history.GetLast(120)
	graph := CreateASCIIGraph(historyValues, graphWidth-2, graphHeight)

	borderStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue)

	// Calculate border sections for clean title placement
	leftBorderLen := (graphWidth - titleWidth - 2) / 2
	rightBorderLen := graphWidth - titleWidth - leftBorderLen - 2
	
	topBorder := config.BorderTopLeft + 
		strings.Repeat(config.BorderHorizontal, leftBorderLen) + 
		title + 
		strings.Repeat(config.BorderHorizontal, rightBorderLen) + 
		config.BorderTopRight

	var result strings.Builder
	result.WriteString(borderStyle.Render(topBorder))
	result.WriteString("\n")

	graphLines := strings.Split(graph, "\n")
	for _, line := range graphLines {
		result.WriteString(borderStyle.Render(config.BorderVertical))
		result.WriteString(line)
		// Pad to width
		currentWidth := lipgloss.Width(line)
		padding := graphWidth - currentWidth - 2
		if padding > 0 {
			result.WriteString(strings.Repeat(" ", padding))
		}
		result.WriteString(borderStyle.Render(config.BorderVertical))
		result.WriteString("\n")
	}

	bottomBorder := config.BorderBottomLeft + 
		strings.Repeat(config.BorderHorizontal, graphWidth-2) + 
		config.BorderBottomRight
	result.WriteString(borderStyle.Render(bottomBorder))

	return result.String()
}

func (m Model) renderMemoryInfo() string {
	barWidth := m.width
	memBar := CreateMemoryBar(m.metrics.MemoryUsed, m.metrics.MemoryTotal, m.metrics.MemoryUsage, barWidth)
	return memBar
}

func (m Model) renderBottomInfo() string {
	loadStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonGreen)

	processStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonPink)

	uptimeStyle := lipgloss.NewStyle().
		Foreground(config.Colors.Yellow)

	uptime := formatDuration(m.metrics.Uptime)
	runtime := formatDuration(time.Since(m.startTime))

	info := fmt.Sprintf(
		"Load: %s  Processes: %s  Uptime: %s  Runtime: %s",
		loadStyle.Render(fmt.Sprintf("%.2f %.2f %.2f", 
			m.metrics.LoadAverage[0], 
			m.metrics.LoadAverage[1], 
			m.metrics.LoadAverage[2])),
		processStyle.Render(fmt.Sprintf("%d", m.metrics.ProcessCount)),
		uptimeStyle.Render(uptime),
		uptimeStyle.Render(runtime),
	)

	separator := lipgloss.NewStyle().
		Foreground(config.Colors.DimGray).
		Render(strings.Repeat("─", m.width))

	return separator + "\n" + info
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm %ds", minutes, int(d.Seconds())%60)
}

func (m Model) renderHelpScreen() string {
	var b strings.Builder
	
	// Header
	titleStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonGreen).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)
	
	b.WriteString(titleStyle.Render("CPU Monitor - Help"))
	b.WriteString("\n")
	b.WriteString(BorderStyle.Render(strings.Repeat("═", m.width)))
	b.WriteString("\n\n")
	
	// Section styles
	sectionStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonPurple).
		Bold(true).
		UnderlineSpaces(false).
		Underline(true)
	
	keyStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue).
		Bold(true).
		Width(15).
		Align(lipgloss.Left)
	
	descStyle := lipgloss.NewStyle().
		Foreground(config.Colors.Foreground)
	
	// Keyboard shortcuts
	b.WriteString(sectionStyle.Render("Keyboard Shortcuts"))
	b.WriteString("\n\n")
	
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"h", "Toggle this help screen"},
		{"q, Ctrl+C", "Quit the application"},
		{"r", "Reset CPU history"},
		{"p", "Pause/unpause monitoring"},
	}
	
	for _, s := range shortcuts {
		b.WriteString("  ")
		b.WriteString(keyStyle.Render(s.key))
		b.WriteString(descStyle.Render(s.desc))
		b.WriteString("\n")
	}
	
	// Display sections
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Display Sections"))
	b.WriteString("\n\n")
	
	sections := []struct {
		name string
		desc string
	}{
		{"Total CPU", "Overall system CPU usage percentage"},
		{"Moving Avg", "10-sample moving average of CPU usage"},
		{"Core Bars", "Individual CPU core usage (multi-column layout for many cores)"},
		{"CPU History", "60-second graph of CPU usage over time"},
		{"Memory", "System RAM usage and availability"},
		{"System Info", "Load average, process count, uptime"},
	}
	
	for _, s := range sections {
		b.WriteString("  ")
		b.WriteString(keyStyle.Render(s.name))
		b.WriteString(descStyle.Render(s.desc))
		b.WriteString("\n")
	}
	
	// Color indicators
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Color Indicators"))
	b.WriteString("\n\n")
	
	colorInfo := []struct {
		color lipgloss.Style
		range_ string
		meaning string
	}{
		{GreenStyle, "0-30%", "Low usage"},
		{BlueStyle, "30-50%", "Light usage"},
		{YellowStyle, "50-70%", "Moderate usage"},
		{OrangeStyle, "70-90%", "High usage"},
		{RedStyle, "90-100%", "Critical usage"},
	}
	
	// Special style for color indicator labels to match other sections
	// Account for the 4-char color block + 1 space = 5 chars offset
	colorKeyStyle := lipgloss.NewStyle().
		Foreground(config.Colors.NeonBlue).
		Bold(true).
		Width(10).  // 15 - 5 (for color block and space)
		Align(lipgloss.Left)
	
	for _, c := range colorInfo {
		b.WriteString("  ")
		b.WriteString(c.color.Render("████"))
		b.WriteString(" ")
		b.WriteString(colorKeyStyle.Render(c.range_))
		b.WriteString(descStyle.Render(c.meaning))
		b.WriteString("\n")
	}
	
	// Command line options
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Command Line Options"))
	b.WriteString("\n\n")
	
	options := []struct {
		flag string
		desc string
	}{
		{"-refresh ms", "Set refresh rate in milliseconds (100-5000, default: 500)"},
		{"-history n", "Number of history points to keep (default: 120)"},
		{"-avg n", "Moving average window size (default: 10)"},
		{"-help", "Show command line help"},
	}
	
	for _, o := range options {
		b.WriteString("  ")
		b.WriteString(keyStyle.Render(o.flag))
		b.WriteString(descStyle.Render(o.desc))
		b.WriteString("\n")
	}
	
	// Footer
	b.WriteString("\n")
	b.WriteString(BorderStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")
	
	footerStyle := lipgloss.NewStyle().
		Foreground(config.Colors.DimGray).
		Align(lipgloss.Center).
		Width(m.width)
	
	b.WriteString(footerStyle.Render("Press 'h' to return to monitoring"))
	
	return b.String()
}
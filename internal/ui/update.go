package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if !m.paused {
			m.collectMetrics()
			m.spinnerFrame++
			m.lastUpdate = time.Time(msg)
		}
		return m, tickCmd(m.config.RefreshRate)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		
		case "r":
			if !m.showHelp {
				m.resetHistory()
			}
			return m, nil
		
		case "p":
			if !m.showHelp {
				m.paused = !m.paused
			}
			return m, nil
		
		case "h":
			m.showHelp = !m.showHelp
			return m, nil
		}
	}

	return m, nil
}
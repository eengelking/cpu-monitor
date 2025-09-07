package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/cpu-monitor/internal/config"
	"github.com/user/cpu-monitor/internal/metrics"
)

type Model struct {
	metrics       *metrics.CPUMetrics
	collector     *metrics.Collector
	history       *metrics.History
	coreHistories []*metrics.History
	config        config.Config
	width         int
	height        int
	paused        bool
	showHelp      bool
	spinnerFrame  int
	lastUpdate    time.Time
	startTime     time.Time
	err           error
}

func NewModel(cfg config.Config) Model {
	collector := metrics.NewCollector()
	initialMetrics, _ := collector.Collect()
	
	var coreHistories []*metrics.History
	if initialMetrics != nil {
		for range initialMetrics.PerCoreUsage {
			coreHistories = append(coreHistories, 
				metrics.NewHistory(cfg.HistorySize, cfg.MovingAvgSize))
		}
	}
	
	return Model{
		metrics:       initialMetrics,
		collector:     collector,
		history:       metrics.NewHistory(cfg.HistorySize, cfg.MovingAvgSize),
		coreHistories: coreHistories,
		config:        cfg,
		width:         80,
		height:        24,
		paused:        false,
		spinnerFrame:  0,
		startTime:     time.Now(),
		lastUpdate:    time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(m.config.RefreshRate),
		tea.WindowSize(),
	)
}

type tickMsg time.Time

func tickCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) collectMetrics() {
	newMetrics, err := m.collector.Collect()
	if err != nil {
		m.err = err
		return
	}
	
	m.metrics = newMetrics
	m.err = nil
	
	if m.metrics.TotalUsage > 0 {
		m.history.Add(m.metrics.TotalUsage)
	}
	
	for i, usage := range m.metrics.PerCoreUsage {
		if i < len(m.coreHistories) {
			m.coreHistories[i].Add(usage)
		}
	}
	
	if len(m.coreHistories) < len(m.metrics.PerCoreUsage) {
		for i := len(m.coreHistories); i < len(m.metrics.PerCoreUsage); i++ {
			hist := metrics.NewHistory(m.config.HistorySize, m.config.MovingAvgSize)
			hist.Add(m.metrics.PerCoreUsage[i])
			m.coreHistories = append(m.coreHistories, hist)
		}
	}
}

func (m *Model) resetHistory() {
	m.history.Reset()
	for _, h := range m.coreHistories {
		h.Reset()
	}
}
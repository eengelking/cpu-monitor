package metrics

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type CPUMetrics struct {
	TotalUsage    float64
	PerCoreUsage  []float64
	Temperature   float64
	Frequency     float64
	ModelName     string
	CoreCount     int
	ThreadCount   int
	ProcessCount  int
	LoadAverage   [3]float64
	MemoryUsage   float64
	MemoryTotal   uint64
	MemoryUsed    uint64
	Uptime        time.Duration
	Timestamp     time.Time
}

type Collector struct {
	lastPerCPU []float64
	// Cached static values
	cpuModelName string
	coreCount    int
	threadCount  int
	// Throttle expensive operations
	lastProcessUpdate time.Time
	processCount      int
	lastTempUpdate    time.Time
	temperature       float64
}

func NewCollector() *Collector {
	c := &Collector{
		threadCount: runtime.NumCPU(),
	}
	
	// Cache static CPU info
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		c.cpuModelName = cpuInfo[0].ModelName
		c.coreCount = int(cpuInfo[0].Cores)
	}
	
	return c
}

func (c *Collector) Collect() (*CPUMetrics, error) {
	metrics := &CPUMetrics{
		Timestamp: time.Now(),
	}

	// Get per-core CPU usage with 100ms sampling interval (single call for efficiency)
	perCorePercent, err := cpu.Percent(100*time.Millisecond, true)
	if err == nil && len(perCorePercent) > 0 {
		metrics.PerCoreUsage = perCorePercent
		
		// Calculate total CPU usage from per-core data (average of all cores)
		var total float64
		for _, coreUsage := range perCorePercent {
			total += coreUsage
		}
		metrics.TotalUsage = total / float64(len(perCorePercent))
	}

	// Use cached static values
	metrics.ModelName = c.cpuModelName
	metrics.CoreCount = c.coreCount
	metrics.ThreadCount = c.threadCount
	
	// Get current frequency (this can change)
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		metrics.Frequency = cpuInfo[0].Mhz
	}

	// Update temperature every 2 seconds
	if time.Since(c.lastTempUpdate) > 2*time.Second {
		temps, err := host.SensorsTemperatures()
		if err == nil {
			for _, temp := range temps {
				if temp.SensorKey == "coretemp_core_0" || temp.SensorKey == "TC0P" || temp.SensorKey == "CPU" {
					c.temperature = temp.Temperature
					break
				}
			}
			if c.temperature == 0 && len(temps) > 0 {
				c.temperature = temps[0].Temperature
			}
		}
		c.lastTempUpdate = time.Now()
	}
	metrics.Temperature = c.temperature

	loadAvg, err := load.Avg()
	if err == nil {
		metrics.LoadAverage = [3]float64{loadAvg.Load1, loadAvg.Load5, loadAvg.Load15}
	}

	// Update process count every 3 seconds
	if time.Since(c.lastProcessUpdate) > 3*time.Second {
		processes, err := process.Processes()
		if err == nil {
			c.processCount = len(processes)
		}
		c.lastProcessUpdate = time.Now()
	}
	metrics.ProcessCount = c.processCount

	vmStat, err := mem.VirtualMemory()
	if err == nil {
		metrics.MemoryUsage = vmStat.UsedPercent
		metrics.MemoryTotal = vmStat.Total
		metrics.MemoryUsed = vmStat.Used
	}

	hostInfo, err := host.Info()
	if err == nil {
		metrics.Uptime = time.Duration(hostInfo.Uptime) * time.Second
	}

	return metrics, nil
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}
package metrics

import (
	"sync"
)

type History struct {
	mu            sync.RWMutex
	values        []float64
	maxSize       int
	movingAvgSize int
	currentIndex  int
}

func NewHistory(maxSize, movingAvgSize int) *History {
	return &History{
		values:        make([]float64, 0, maxSize),
		maxSize:       maxSize,
		movingAvgSize: movingAvgSize,
	}
}

func (h *History) Add(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.values) < h.maxSize {
		h.values = append(h.values, value)
	} else {
		h.values[h.currentIndex] = value
		h.currentIndex = (h.currentIndex + 1) % h.maxSize
	}
}

func (h *History) GetValues() []float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]float64, len(h.values))
	if len(h.values) < h.maxSize {
		copy(result, h.values)
	} else {
		for i := 0; i < len(h.values); i++ {
			idx := (h.currentIndex + i) % h.maxSize
			result[i] = h.values[idx]
		}
	}
	return result
}

func (h *History) GetMovingAverage() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) == 0 {
		return 0
	}

	count := h.movingAvgSize
	if len(h.values) < count {
		count = len(h.values)
	}

	sum := 0.0
	if len(h.values) < h.maxSize {
		start := len(h.values) - count
		if start < 0 {
			start = 0
		}
		for i := start; i < len(h.values); i++ {
			sum += h.values[i]
		}
	} else {
		for i := 0; i < count; i++ {
			idx := (h.currentIndex - count + i + h.maxSize) % h.maxSize
			sum += h.values[idx]
		}
	}

	return sum / float64(count)
}

func (h *History) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.values = make([]float64, 0, h.maxSize)
	h.currentIndex = 0
}

func (h *History) GetLast(n int) []float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.values) == 0 {
		return []float64{}
	}

	if n > len(h.values) {
		n = len(h.values)
	}

	result := make([]float64, n)
	if len(h.values) < h.maxSize {
		start := len(h.values) - n
		if start < 0 {
			start = 0
		}
		copy(result, h.values[start:])
	} else {
		for i := 0; i < n; i++ {
			idx := (h.currentIndex - n + i + h.maxSize) % h.maxSize
			result[i] = h.values[idx]
		}
	}

	return result
}
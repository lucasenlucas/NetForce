package metrics

import (
	"sort"
	"sync"
	"time"
)

// Result holds data for a single completed request
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      string
	Timestamp  time.Time
}

// Collector stores and aggregates all results thread-safely
type Collector struct {
	mu        sync.Mutex
	results   []Result
	startTime time.Time
	endTime   time.Time
	done      bool
}

func NewCollector() *Collector {
	return &Collector{results: make([]Result, 0, 2048)}
}

func (c *Collector) Start() {
	c.mu.Lock()
	c.startTime = time.Now()
	c.mu.Unlock()
}

func (c *Collector) Stop() {
	c.mu.Lock()
	c.endTime = time.Now()
	c.done = true
	c.mu.Unlock()
}

func (c *Collector) Add(r Result) {
	c.mu.Lock()
	c.results = append(c.results, r)
	c.mu.Unlock()
}

func (c *Collector) Snapshot() []Result {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Result, len(c.results))
	copy(out, c.results)
	return out
}

func (c *Collector) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.results)
}

func (c *Collector) Successes() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, r := range c.results {
		if r.StatusCode >= 200 && r.StatusCode < 400 {
			n++
		}
	}
	return n
}

func (c *Collector) Errors() int {
	return c.Total() - c.Successes()
}

func (c *Collector) AvgDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.results) == 0 {
		return 0
	}
	var sum time.Duration
	for _, r := range c.results {
		sum += r.Duration
	}
	return sum / time.Duration(len(c.results))
}

func (c *Collector) MaxDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	var max time.Duration
	for _, r := range c.results {
		if r.Duration > max {
			max = r.Duration
		}
	}
	return max
}

func (c *Collector) SuccessRate() float64 {
	t := c.Total()
	if t == 0 {
		return 0
	}
	return float64(c.Successes()) / float64(t) * 100
}

func (c *Collector) TestDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done {
		return c.endTime.Sub(c.startTime)
	}
	return time.Since(c.startTime)
}

func (c *Collector) Latencies() []time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.results) == 0 {
		return nil
	}
	out := make([]time.Duration, len(c.results))
	for i, r := range c.results {
		out[i] = r.Duration
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

func (c *Collector) Percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)) * p)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func (c *Collector) StatusCodes() map[int]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	counts := make(map[int]int)
	for _, r := range c.results {
		counts[r.StatusCode]++
	}
	return counts
}

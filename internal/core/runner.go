package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lucasenlucas/netforce/internal/metrics"
)

// RunConfig holds all parameters needed to execute a test
type RunConfig struct {
	URL      string
	Rate     int
	Threads  int
	Duration int
	ModeName string
	Timeout  int
	Live     bool
}

// Run launches the test and returns the collected metrics
func Run(ctx context.Context, cfg RunConfig) *metrics.Collector {
	col := metrics.NewCollector()
	col.Start()

	client := buildClient(cfg.Threads, cfg.Timeout)

	var pacer Pacer
	switch cfg.ModeName {
	case "ramp":
		pacer = NewRampPacer(cfg.Rate, cfg.Threads, cfg.Duration)
	case "spike":
		pacer = NewSpikePacer(cfg.Rate, cfg.Threads, cfg.Duration)
	case "pulse":
		pacer = NewPulsePacer(cfg.Rate, cfg.Threads)
	case "unlimited":
		pacer = NewUnlimitedPacer()
	default:
		pacer = NewConstantPacer(cfg.Rate, cfg.Threads)
	}

	deadline := time.Now().Add(time.Duration(cfg.Duration) * time.Second)
	runCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < cfg.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(runCtx, client, cfg.URL, pacer, col)
		}()
	}

	if cfg.Live {
		wg.Add(1)
		go func() {
			defer wg.Done()
			liveDashboard(runCtx, col)
		}()
	}

	wg.Wait()
	col.Stop()
	return col
}

func liveDashboard(ctx context.Context, col *metrics.Collector) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// safely clear line
			fmt.Print("\r\033[2K\r")
			return
		case <-ticker.C:
			reqs := col.Total()
			var rps float64
			dur := col.TestDuration().Seconds()
			if dur > 0 {
				rps = float64(reqs) / dur
			}
			succRate := col.SuccessRate()
			avgLat := col.AvgDuration().Milliseconds()
			fmt.Printf("\r  Live Stats: %d reqs | %.0f req/s | %.1f%% success | %dms avg latency ...", reqs, rps, succRate, avgLat)
		}
	}
}

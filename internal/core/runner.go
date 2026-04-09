package core

import (
	"context"
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
}

// Run launches the test and returns the collected metrics
func Run(ctx context.Context, cfg RunConfig) *metrics.Collector {
	col := metrics.NewCollector()
	col.Start()

	client := buildClient(cfg.Threads, cfg.Timeout)

	var pacer Pacer
	switch cfg.ModeName {
	// TODO: add ramp, spike, pulse cases once their pacers are implemented
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

	wg.Wait()
	col.Stop()
	return col
}

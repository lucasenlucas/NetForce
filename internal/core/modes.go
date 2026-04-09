package core

import (
	"context"
	"time"
)

// Pacer controls the rate of requests each worker sends
type Pacer interface {
	Wait(ctx context.Context) error
}

// ConstantPacer spaces requests evenly to achieve the target rate per worker
type ConstantPacer struct {
	delay time.Duration
}

// NewConstantPacer creates a pacer that divides the target rate across workers
func NewConstantPacer(ratePerSecond, threads int) *ConstantPacer {
	var delay time.Duration
	if ratePerSecond > 0 && threads > 0 {
		// Each worker is responsible for (rate / threads) requests per second
		perWorker := float64(ratePerSecond) / float64(threads)
		delay = time.Duration(float64(time.Second) / perWorker)
	}
	return &ConstantPacer{delay: delay}
}

func (p *ConstantPacer) Wait(ctx context.Context) error {
	if p.delay == 0 {
		return nil
	}
	t := time.NewTimer(p.delay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// TODO: Implement RampPacer  — gradually decrease delay over test duration
// TODO: Implement SpikePacer — sudden burst then return to baseline
// TODO: Implement PulsePacer — alternate between high and low rates

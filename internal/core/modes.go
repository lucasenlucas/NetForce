package core

import (
	"context"
	"math"
	"sync/atomic"
	"time"
)

// Pacer controls the rate of requests each worker sends
type Pacer interface {
	Wait(ctx context.Context) error
}

// ConstantPacer spaces requests evenly at a steady rate
type ConstantPacer struct {
	delay time.Duration
}

func NewConstantPacer(ratePerSecond, threads int) *ConstantPacer {
	var delay time.Duration
	if ratePerSecond > 0 && threads > 0 {
		perWorker := float64(ratePerSecond) / float64(threads)
		delay = time.Duration(float64(time.Second) / perWorker)
	}
	return &ConstantPacer{delay: delay}
}

func (p *ConstantPacer) Wait(ctx context.Context) error {
	return sleepCtx(ctx, p.delay)
}

// RampPacer gradually increases the rate from 10% of max to 100% over the test duration
type RampPacer struct {
	threads   int
	maxRate   int
	startTime time.Time
	duration  time.Duration
}

func NewRampPacer(ratePerSecond, threads, durationSec int) *RampPacer {
	return &RampPacer{
		threads:   threads,
		maxRate:   ratePerSecond,
		startTime: time.Now(),
		duration:  time.Duration(durationSec) * time.Second,
	}
}

func (p *RampPacer) Wait(ctx context.Context) error {
	// Calculate how far through the test we are (0.0 → 1.0)
	elapsed := time.Since(p.startTime)
	progress := math.Min(elapsed.Seconds()/p.duration.Seconds(), 1.0)

	// Current rate ramps from 10% to 100% of max rate
	minRate := math.Max(float64(p.maxRate)*0.10, 1)
	currentRate := minRate + (float64(p.maxRate)-minRate)*progress

	perWorker := currentRate / float64(p.threads)
	delay := time.Duration(float64(time.Second) / perWorker)
	return sleepCtx(ctx, delay)
}

// SpikePacer sends a sudden high burst for the first 20% of the test, then stops each worker
type SpikePacer struct {
	threads   int
	maxRate   int
	startTime time.Time
	duration  time.Duration
	done      atomic.Bool
}

func NewSpikePacer(ratePerSecond, threads, durationSec int) *SpikePacer {
	return &SpikePacer{
		threads:   threads,
		maxRate:   ratePerSecond,
		startTime: time.Now(),
		duration:  time.Duration(durationSec) * time.Second,
	}
}

func (p *SpikePacer) Wait(ctx context.Context) error {
	elapsed := time.Since(p.startTime)
	spikeDuration := time.Duration(float64(p.duration) * 0.20) // burst = first 20%

	if elapsed > spikeDuration {
		// After spike: long sleep so workers idle but don't exit (context will cancel them)
		return sleepCtx(ctx, 10*time.Second)
	}

	// During spike: max rate
	perWorker := float64(p.maxRate) / float64(p.threads)
	delay := time.Duration(float64(time.Second) / perWorker)
	return sleepCtx(ctx, delay)
}

// PulsePacer alternates between high-rate bursts and rest periods
// Pattern: 3s burst at full rate → 2s pause → repeat
type PulsePacer struct {
	threads int
	maxRate int
}

func NewPulsePacer(ratePerSecond, threads int) *PulsePacer {
	return &PulsePacer{threads: threads, maxRate: ratePerSecond}
}

func (p *PulsePacer) Wait(ctx context.Context) error {
	cycleLength := 5 * time.Second  // full cycle = 5s
	burstLength := 3 * time.Second  // send for 3s

	// Where are we in the current cycle?
	posInCycle := time.Duration(time.Now().UnixNano()) % cycleLength

	if posInCycle < burstLength {
		// Burst phase: full speed
		perWorker := float64(p.maxRate) / float64(p.threads)
		delay := time.Duration(float64(time.Second) / perWorker)
		return sleepCtx(ctx, delay)
	}

	// Rest phase: sleep until next burst
	restRemaining := cycleLength - posInCycle
	return sleepCtx(ctx, restRemaining)
}

// sleepCtx sleeps for d duration but returns early if ctx is cancelled
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

package analyzer

import (
	"fmt"
	"github.com/lucasenlucas/netforce/internal/metrics"
)

// AnalyzePerformance compares the first and second half of successful requests
// to detect response time degradation under sustained load
func AnalyzePerformance(results []metrics.Result) string {
	var successes []metrics.Result
	for _, r := range results {
		if r.StatusCode >= 200 && r.StatusCode < 400 {
			successes = append(successes, r)
		}
	}

	if len(successes) < 10 {
		return "Not enough successful responses to analyze performance degradation."
	}

	mid := len(successes) / 2
	firstAvg  := avgDuration(successes[:mid])
	secondAvg := avgDuration(successes[mid:])

	if secondAvg > firstAvg*1.5 {
		pct := (secondAvg/firstAvg)*100 - 100
		return fmt.Sprintf(
			"⚠  Performance Degradation Detected: avg response time grew from %.0fms to %.0fms (+%.0f%%)",
			firstAvg, secondAvg, pct,
		)
	}

	return fmt.Sprintf(
		"✓  Performance stable (first half avg: %.0fms, second half avg: %.0fms)",
		firstAvg, secondAvg,
	)
}

func avgDuration(results []metrics.Result) float64 {
	if len(results) == 0 {
		return 0
	}
	var total float64
	for _, r := range results {
		total += float64(r.Duration.Milliseconds())
	}
	return total / float64(len(results))
}

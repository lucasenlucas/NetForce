package report

import (
	"fmt"
	"os"
	"time"

	"github.com/lucasenlucas/netforce/internal/output"
)

// Generate saves a human-readable test report to a timestamped text file
func Generate(s *output.Summary) (string, error) {
	filename := fmt.Sprintf("netforce_report_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	f, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("could not create report file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "NetForce Test Report\n")
	fmt.Fprintf(f, "Generated : %s\n\n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(f, "Target          : %s\n", s.Target)
	fmt.Fprintf(f, "Feature         : %s\n", s.Feature)
	fmt.Fprintf(f, "Mode            : %s\n", s.Mode)
	fmt.Fprintf(f, "Configured Rate : %d req/s\n\n", s.ConfiguredRate)
	fmt.Fprintf(f, "Duration        : %.1fs\n", s.DurationSec)
	fmt.Fprintf(f, "Total Requests  : %d\n", s.TotalRequests)
	fmt.Fprintf(f, "Successes       : %d\n", s.Successes)
	fmt.Fprintf(f, "Errors          : %d\n", s.Errors)
	fmt.Fprintf(f, "Success Rate    : %.2f%%\n", s.SuccessRate)
	fmt.Fprintf(f, "Error Rate      : %.2f%%\n\n", s.ErrorRate)
	fmt.Fprintf(f, "Avg Latency     : %.0fms\n", s.AvgLatencyMs)
	fmt.Fprintf(f, "Max Latency     : %.0fms\n\n", s.MaxLatencyMs)
	fmt.Fprintf(f, "Rate Limit Hit  : %v\n", s.RateLimitHit)
	if s.Analysis != "" {
		fmt.Fprintf(f, "Analysis        : %s\n", s.Analysis)
	}

	return filename, nil
}

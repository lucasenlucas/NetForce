package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lucasenlucas/netforce/internal/metrics"
)

type Summary struct {
	Target        string  `json:"target"`
	Feature       string  `json:"feature"`
	Mode          string  `json:"mode"`
	ConfiguredRate int    `json:"configured_rate_rps"`
	TotalRequests  int    `json:"total_requests"`
	Successes      int    `json:"successes"`
	Errors         int    `json:"errors"`
	SuccessRate    float64 `json:"success_rate_pct"`
	ErrorRate      float64 `json:"error_rate_pct"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	MaxLatencyMs   float64 `json:"max_latency_ms"`
	DurationSec    float64 `json:"duration_sec"`
	RateLimitHit   bool    `json:"rate_limit_detected"`
	Analysis       string  `json:"performance_analysis,omitempty"`
}

func BuildSummary(target, feature, mode string, rate int, col *metrics.Collector, rateLimitHit bool, analysis string) *Summary {
	return &Summary{
		Target:         target,
		Feature:        feature,
		Mode:           mode,
		ConfiguredRate: rate,
		TotalRequests:  col.Total(),
		Successes:      col.Successes(),
		Errors:         col.Errors(),
		SuccessRate:    col.SuccessRate(),
		ErrorRate:      100 - col.SuccessRate(),
		AvgLatencyMs:   float64(col.AvgDuration().Milliseconds()),
		MaxLatencyMs:   float64(col.MaxDuration().Milliseconds()),
		DurationSec:    col.TestDuration().Seconds(),
		RateLimitHit:   rateLimitHit,
		Analysis:       analysis,
	}
}

func PrintSimple(s *Summary) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════╗")
	color.Cyan("║       NetForce — Test Results            ║")
	color.Cyan("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  %-22s %s\n", "Target:", s.Target)
	fmt.Printf("  %-22s %s  (mode: %s)\n", "Feature:", s.Feature, s.Mode)
	fmt.Printf("  %-22s %d req/s\n", "Configured Rate:", s.ConfiguredRate)
	fmt.Printf("  %-22s %.1fs\n", "Test Duration:", s.DurationSec)
	fmt.Println()
	fmt.Printf("  %-22s %d\n", "Total Requests:", s.TotalRequests)
	color.Green("  %-22s %d\n", "Successes:", s.Successes)
	if s.Errors > 0 {
		color.Red("  %-22s %d\n", "Errors:", s.Errors)
	} else {
		fmt.Printf("  %-22s 0\n", "Errors:")
	}
	fmt.Printf("  %-22s %.2f%%\n", "Success Rate:", s.SuccessRate)
	fmt.Printf("  %-22s %.2f%%\n", "Error Rate:", s.ErrorRate)
	fmt.Println()
	fmt.Printf("  %-22s %.0fms\n", "Avg Latency:", s.AvgLatencyMs)
	fmt.Printf("  %-22s %.0fms\n", "Max Latency:", s.MaxLatencyMs)
	fmt.Println()

	printAnalysisSection(s)
}

func PrintDetailed(s *Summary) {
	PrintSimple(s)
	// TODO: Add per-status-code breakdown and percentile latencies (p50/p90/p99)
	fmt.Println("  [detailed mode] — extended breakdown coming soon")
}

func PrintJSON(s *Summary) {
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))
}

func printAnalysisSection(s *Summary) {
	color.Cyan("  ─────────────────────────────────────────")
	if s.RateLimitHit {
		color.Yellow("  ⚠  Rate Limiting Detected (HTTP 429 observed)")
		fmt.Println("     The server appears to be throttling requests.")
	} else {
		color.Green("  ✓  No rate limiting detected")
	}

	if s.Analysis != "" {
		prefix := "  "
		fmt.Println(prefix + strings.ReplaceAll(s.Analysis, "\n", "\n"+prefix))
	}
	fmt.Println()
}

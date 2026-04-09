package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lucasenlucas/netforce/internal/metrics"
)

type Summary struct {
	Target        string      `json:"target"`
	Feature       string      `json:"feature"`
	Mode          string      `json:"mode"`
	ConfiguredRate int        `json:"configured_rate_rps"`
	TotalRequests  int        `json:"total_requests"`
	Successes      int        `json:"successes"`
	Errors         int        `json:"errors"`
	SuccessRate    float64    `json:"success_rate_pct"`
	ErrorRate      float64    `json:"error_rate_pct"`
	AvgLatencyMs   float64    `json:"avg_latency_ms"`
	MaxLatencyMs   float64    `json:"max_latency_ms"`
	P50LatencyMs   float64    `json:"p50_latency_ms,omitempty"`
	P90LatencyMs   float64    `json:"p90_latency_ms,omitempty"`
	P99LatencyMs   float64    `json:"p99_latency_ms,omitempty"`
	StatusCodes    map[int]int `json:"status_codes,omitempty"`
	DurationSec    float64    `json:"duration_sec"`
	RateLimitHit   bool       `json:"rate_limit_detected"`
	Analysis       string     `json:"performance_analysis,omitempty"`
}

func BuildSummary(target, feature, mode string, rate int, col *metrics.Collector, rateLimitHit bool, analysis string) *Summary {
	lats := col.Latencies()
	p50 := col.Percentile(lats, 0.50)
	p90 := col.Percentile(lats, 0.90)
	p99 := col.Percentile(lats, 0.99)

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
		P50LatencyMs:   float64(p50.Milliseconds()),
		P90LatencyMs:   float64(p90.Milliseconds()),
		P99LatencyMs:   float64(p99.Milliseconds()),
		StatusCodes:    col.StatusCodes(),
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

	color.Cyan("  ─────────────────────────────────────────")
	color.Yellow("  Detailed Latency Breakdown:")
	fmt.Printf("  %-22s %.0fms\n", "p50 (Median):", s.P50LatencyMs)
	fmt.Printf("  %-22s %.0fms\n", "p90:", s.P90LatencyMs)
	fmt.Printf("  %-22s %.0fms\n", "p99:", s.P99LatencyMs)
	fmt.Println()

	color.Yellow("  Status Code Distribution:")
	if len(s.StatusCodes) == 0 {
		fmt.Println("  No responses received.")
	} else {
		for code, count := range s.StatusCodes {
			label := fmt.Sprintf("[%d]", code)
			if code >= 200 && code < 300 {
				color.Green("  %-22s %d requests\n", label, count)
			} else if code >= 400 && code < 500 {
				color.Yellow("  %-22s %d requests\n", label, count)
			} else {
				color.Red("  %-22s %d requests\n", label, count)
			}
		}
	}
	fmt.Println()
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

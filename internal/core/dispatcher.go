package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lucasenlucas/netforce/internal/analyzer"
	"github.com/lucasenlucas/netforce/internal/cli"
	"github.com/lucasenlucas/netforce/internal/output"
	"github.com/lucasenlucas/netforce/internal/report"
	"github.com/lucasenlucas/netforce/internal/validate"
)

const (
	safeMaxRate     = 20
	safeMaxThreads  = 5
	safeMaxDuration = 15
)

// Dispatch reads the -f feature flag and routes to the correct handler
func Dispatch(cfg *cli.Config) {
	switch strings.ToLower(cfg.Feature) {
	case "explain":
		runExplain()
	case "benchmark":
		RunBenchmark()
	case "quick":
		runQuick(cfg)
	case "stress":
		runTest(cfg, "stress", "constant")
	case "ramp":
		runTest(cfg, "ramp", "ramp")
	case "spike":
		runTest(cfg, "spike", "spike")
	case "pulse":
		runTest(cfg, "pulse", "pulse")
	default:
		color.Red("Unknown feature: %q\nRun with --help to see available features.", cfg.Feature)
		os.Exit(1)
	}
}

func runExplain() {
	color.Cyan(`
 ███╗   ██╗███████╗████████╗███████╗ ██████╗ ██████╗  ██████╗███████╗
 ████╗  ██║██╔════╝╚══██╔══╝██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔════╝
 ██╔██╗ ██║█████╗     ██║   █████╗  ██║   ██║██████╔╝██║     █████╗
 ██║╚██╗██║██╔══╝     ██║   ██╔══╝  ██║   ██║██╔══██╗██║     ██╔══╝
 ██║ ╚████║███████╗   ██║   ██║     ╚██████╔╝██║  ██║╚██████╗███████╗
 ╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝╚══════╝
`)
	color.White("  A tool made by Lucas Mangroelal — part of the NET Toolkit")
	fmt.Println()
	color.Cyan("  What is NetForce?")

	fmt.Println("  NetForce is a performance and resilience testing tool for websites.")
	fmt.Println("  It simulates real visitor traffic so you can see how your server")
	fmt.Println("  behaves when many people use it at the same time.")
	fmt.Println()
	color.Yellow("  Features / Modes:")
	fmt.Println("  -f stress   Sends a constant, steady stream of requests.")
	fmt.Println("              Good for a baseline: how does the server handle X users at once?")
	fmt.Println()
	fmt.Println("  -f ramp     Gradually increases traffic from low to high.")
	fmt.Println("              Useful for finding the point where the server starts to struggle.")
	fmt.Println()
	fmt.Println("  -f spike    Sends a sudden large burst of traffic for a short period.")
	fmt.Println("              Simulates a flash sale or a viral post sending many visitors at once.")
	fmt.Println()
	fmt.Println("  -f pulse    Sends traffic in repeated waves (burst, pause, burst, pause).")
	fmt.Println("              Models real-world traffic patterns with peaks and quiet periods.")
	fmt.Println()
	fmt.Println("  -f quick    Beginner-friendly test with very safe default settings.")
	fmt.Println("              Great for a quick first look at how a server responds.")
	fmt.Println()
	color.Red("  IMPORTANT: Only test systems you own or have explicit permission to test.")
	fmt.Println()
}

func runQuick(cfg *cli.Config) {
	if !confirmPermission(cfg.Domain) {
		return
	}

	color.Green("\n  Running QUICK mode (safe defaults)")
	fmt.Println("  Rate: 10 req/s | Threads: 2 | Duration: 10s")
	fmt.Println()

	url := buildURL(cfg.Domain, "/", false)
	runCfg := RunConfig{URL: url, Rate: 10, Threads: 2, Duration: 10, ModeName: "constant", Timeout: 10, Live: cfg.Live}
	col := Run(context.Background(), runCfg)

	results := col.Snapshot()
	rateLimitHit := analyzer.DetectRateLimit(results)
	analysis := ""

	sum := output.BuildSummary(cfg.Domain, "quick", "constant", 10, col, rateLimitHit, analysis)
	output.PrintSimple(sum)
}

func runTest(cfg *cli.Config, feature, defaultMode string) {
	// Validate flags
	if err := validate.Rate(cfg.Rate); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
	if err := validate.Threads(cfg.Threads); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
	if err := validate.Duration(cfg.Duration); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
	if err := validate.Output(cfg.Output); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	// Mode resolution: explicit --mode flag overrides feature default
	mode := defaultMode
	if cfg.Mode != "" {
		if err := validate.Mode(cfg.Mode); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		mode = cfg.Mode
	}

	// Apply safe mode caps
	rate, threads, duration := cfg.Rate, cfg.Threads, cfg.Duration
	if cfg.Safe {
		if rate > safeMaxRate       { rate = safeMaxRate }
		if threads > safeMaxThreads { threads = safeMaxThreads }
		if duration > safeMaxDuration { duration = safeMaxDuration }
		color.Yellow("  Safe mode active — limits capped (rate: %d, threads: %d, duration: %ds)", rate, threads, duration)
	}

	if !confirmPermission(cfg.Domain) {
		return
	}

	url := buildURL(cfg.Domain, cfg.Path, cfg.ForceHTTPS)

	color.Cyan("\n  Starting NetForce test...")
	fmt.Printf("  Target: %s\n", url)
	fmt.Printf("  Feature: %s | Mode: %s | Rate: %d req/s | Threads: %d | Duration: %ds\n\n",
		feature, mode, rate, threads, duration)

	runCfg := RunConfig{URL: url, Rate: rate, Threads: threads, Duration: duration, ModeName: mode, Timeout: cfg.Timeout, Live: cfg.Live}
	col := Run(context.Background(), runCfg)

	results := col.Snapshot()

	rateLimitHit := false
	if cfg.DetectRateLimit {
		rateLimitHit = analyzer.DetectRateLimit(results)
	}

	analysis := ""
	if cfg.AnalyzePerf {
		analysis = analyzer.AnalyzePerformance(results)
	}

	sum := output.BuildSummary(cfg.Domain, feature, mode, rate, col, rateLimitHit, analysis)

	switch strings.ToLower(cfg.Output) {
	case "json":
		output.PrintJSON(sum)
	case "detailed":
		output.PrintDetailed(sum)
	default:
		output.PrintSimple(sum)
	}

	if cfg.Report {
		filename, err := report.Generate(sum)
		if err != nil {
			color.Red("  Could not save report: %v", err)
		} else {
			color.Green("  Report saved → %s", filename)
		}
	}
}

func confirmPermission(domain string) bool {
	color.Yellow("\n  ⚠  You are about to run a load test against: %s", domain)
	color.Yellow("  Do you have authorized permission to test this target? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp == "y" || resp == "yes" {
		return true
	}
	color.Red("  Test cancelled. Authorized permission is required.")
	return false
}

func buildURL(domain, path string, forceHTTPS bool) string {
	scheme := "http"
	if forceHTTPS || strings.HasPrefix(domain, "https://") {
		scheme = "https"
	}
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf("%s://%s%s", scheme, domain, path)
}

package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Config holds all parsed CLI flags for NetForce
type Config struct {
	Domain            string
	Feature           string
	Rate              int
	Threads           int
	Duration          int
	Mode              string
	Path              string
	Timeout           int
	ForceHTTPS        bool
	Live              bool
	Output            string
	Report            bool
	Safe              bool
	Confirm           bool
	DetectRateLimit   bool
	AnalyzePerf       bool
}

func Parse() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Domain, "d", "", "Target domain or URL (e.g. example.com)")
	flag.StringVar(&cfg.Domain, "domain", "", "Target domain or URL (e.g. example.com)")

	flag.StringVar(&cfg.Feature, "f", "", "Feature to run: stress, ramp, spike, pulse, quick, explain")
	flag.StringVar(&cfg.Feature, "feature", "", "Feature to run: stress, ramp, spike, pulse, quick, explain")

	flag.IntVar(&cfg.Rate, "r", 10, "Requests per second to send")
	flag.IntVar(&cfg.Rate, "rate", 10, "Requests per second to send")

	flag.IntVar(&cfg.Threads, "t", 5, "Number of concurrent workers")
	flag.IntVar(&cfg.Threads, "threads", 5, "Number of concurrent workers")

	flag.IntVar(&cfg.Duration, "duration", 10, "How long to run the test in seconds")
	flag.IntVar(&cfg.Timeout, "timeout", 10, "Max wait time per request in seconds")

	flag.StringVar(&cfg.Mode, "mode", "", "Override load mode: constant, ramp, spike, pulse")
	flag.StringVar(&cfg.Path, "path", "/", "URL path to target (e.g. /api/health)")
	flag.StringVar(&cfg.Output, "output", "simple", "Output format: simple, detailed, json")

	flag.BoolVar(&cfg.ForceHTTPS, "https", false, "Force HTTPS even if domain has no scheme")
	flag.BoolVar(&cfg.Live, "live", false, "Show live stats while the test runs")
	flag.BoolVar(&cfg.Report, "report", false, "Save a text report after the test")
	flag.BoolVar(&cfg.Safe, "safe", false, "Enable safe mode (caps rate, threads, duration)")
	flag.BoolVar(&cfg.Confirm, "confirm", false, "Ask for confirmation before starting the test")
	flag.BoolVar(&cfg.DetectRateLimit, "detect-rate-limit", false, "Detect if the target is rate limiting responses")
	flag.BoolVar(&cfg.AnalyzePerf, "analyze-performance", false, "Analyze response time degradation during the test")

	flag.Usage = PrintUsage
	flag.Parse()

	return cfg
}

func PrintUsage() {
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
	fmt.Println("  Performance & Resilience Testing Tool")
	fmt.Println("  For authorized testing on systems you own or have explicit permission to test.")
	fmt.Println()
	color.Yellow("  DISCLAIMER: Unauthorized use is strictly prohibited.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  netforce -d <domain> -f <feature> [options]")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  netforce -d example.com -f quick")
	fmt.Println("  netforce -d example.com -f stress -r 50 -t 5 --duration 30")
	fmt.Println("  netforce -d example.com -f ramp -r 100 --duration 60 --report")
	fmt.Println("  netforce -f explain")
	fmt.Println()
	fmt.Println("PRIMARY FLAGS:")
	fmt.Println("  -d, --domain    Target domain or URL")
	fmt.Println("  -f, --feature   Feature: stress | ramp | spike | pulse | quick | explain")
	fmt.Println("  -r, --rate      Requests per second (default: 10)")
	fmt.Println("  -t, --threads   Concurrent workers (default: 5)")
	fmt.Println("      --duration  Test duration in seconds (default: 10)")
	fmt.Println("      --timeout   Request timeout in seconds (default: 10)")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("      --mode              Override mode: constant | ramp | spike | pulse")
	fmt.Println("      --path              URL path to target (default: /)")
	fmt.Println("      --https             Force HTTPS")
	fmt.Println("      --live              Show live stats during test")
	fmt.Println("      --output            Output format: simple | detailed | json")
	fmt.Println("      --report            Save report to file after test")
	fmt.Println("      --safe              Enable safe mode caps")
	fmt.Println("      --confirm           Ask for confirmation before starting")
	fmt.Println("      --detect-rate-limit Detect HTTP 429 rate limiting")
	fmt.Println("      --analyze-performance Analyze response time degradation")
	fmt.Println()

	// Make sure we exit cleanly
	os.Exit(0)
}

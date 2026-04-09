package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"

	"github.com/fatih/color"
)

func RunBenchmark() {
	color.Cyan("\n  Starting NetForce Hardware Calibration & Benchmark")
	color.White("  This will measure the absolute maximum requests per second (RPS)")
	color.White("  your current CPU and local network stack can generate.")
	fmt.Println()

	color.Yellow("  Warming up local test server...")

	// 1. Create a lightning-fast local dummy server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// 2. Configure an extreme load configuration
	// We use 100 threads per CPU core to saturate connection dialing capability
	threads := runtime.NumCPU() * 100
	if threads < 200 {
		threads = 200 // minimum baseline
	}
	duration := 5 // 5 seconds is enough to find the ceiling without burning the laptop

	color.Yellow("  Saturating local cores (%d threads) for %d seconds...", threads, duration)
	fmt.Println()

	runCfg := RunConfig{
		URL:      ts.URL,
		Rate:     0, // Unused by unlimited pacer
		Threads:  threads,
		Duration: duration,
		ModeName: "unlimited",
		Timeout:  5,
		Live:     true,
	}

	// 3. Blast the local endpoint
	col := Run(context.Background(), runCfg)

	// 4. Calculate performance
	totalReqs := col.Total()
	actualDur := col.TestDuration().Seconds()
	
	var rps int
	if actualDur > 0 {
		rps = int(float64(totalReqs) / actualDur)
	}

	// 5. Provide analysis and recommendation
	fmt.Println()
	color.Green("  Benchmark Complete!")
	fmt.Printf("  Max Generated Load: %d req/s\n", rps)
	fmt.Printf("  Total Requests sent in %.1fs: %d\n", actualDur, totalReqs)
	fmt.Println()

	// Advise a safe "extreme" limit using 80% to leave headroom for OS/Background tasks during real tests
	recommendedMaxRate := int(float64(rps) * 0.8)
	// Round down to nearest 500 for a clean number
	if recommendedMaxRate > 500 {
		recommendedMaxRate = (recommendedMaxRate / 500) * 500
	} else if recommendedMaxRate == 0 {
		recommendedMaxRate = 100
	}

	color.Cyan("  💡 Configuration Advice for Extreme Tests:")
	fmt.Println("  To push an external server to YOUR absolute limit without crashing your own PC, use:")
	color.White("  --rate %d --threads %d\n", recommendedMaxRate, threads)
	fmt.Println()
}

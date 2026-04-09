package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

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
	
	// 2. We explicitly start an escalation loop to find the exact local bottleneck
	var maxRps int
	var maxThreads int

	currentThreads := 500
	step := 500

	for {
		color.Yellow("  Testing capacity at %d threads... (2 second burst)", currentThreads)

		runCfg := RunConfig{
			URL:      ts.URL,
			Rate:     0, // Unused by unlimited pacer
			Threads:  currentThreads,
			Duration: 2, // Short bursts to find limits quickly
			ModeName: "unlimited",
			Timeout:  5,
			Live:     false, // Disable interactive UI spam during the loop
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
		successRate := col.SuccessRate()

		fmt.Printf("\r\033[2K\r  -> %d threads | %.1f%% success | %d req/s \n", currentThreads, successRate, rps)

		// Limit Case A: Socket Exhaustion / Network Crash
		if successRate < 98.0 {
			color.Red("  [LIMIT REACHED] System connection limit or socket exhaustion detected!")
			break
		}

		// Limit Case B: CPU Plateau (RPS dropped or stopped growing)
		if maxRps > 0 && rps <= int(float64(maxRps)*1.02) {
			color.Yellow("  [LIMIT REACHED] Hardware throughput peaked and plateaud.")
			break
		}

		// Save the new peak and escalate
		maxRps = rps
		maxThreads = currentThreads

		currentThreads += step
		if currentThreads > 50000 {
			color.Yellow("  [CAPPED] Reached extreme ceiling, stopping for safety.")
			break
		}
	}

	// 5. Provide analysis and recommendation
	fmt.Println()
	color.Green("  Benchmark Calibration Complete!")
	fmt.Printf("  Max Safe Threads: %d\n", maxThreads)
	fmt.Printf("  Max Generated Load: %d req/s\n", maxRps)
	fmt.Println()

	// Advise a safe "extreme" limit using 80% to leave headroom for OS/Background tasks during real tests
	recommendedMaxRate := int(float64(maxRps) * 0.8)
	// Round down to nearest 500 for a clean number
	if recommendedMaxRate > 500 {
		recommendedMaxRate = (recommendedMaxRate / 500) * 500
	} else if recommendedMaxRate == 0 {
		recommendedMaxRate = 100
	}

	color.Cyan("  💡 Configuration Advice for Extreme Tests:")
	fmt.Println("  To push an external server to YOUR absolute limit without crashing your own PC, use:")
	color.White("  --rate %d --threads %d\n", recommendedMaxRate, maxThreads)
	fmt.Println()
}

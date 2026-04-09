package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/lucasenlucas/netforce/internal/cli"
	"github.com/lucasenlucas/netforce/internal/core"
	"github.com/lucasenlucas/netforce/internal/validate"
)

func main() {
	cfg := cli.Parse()

	// Feature is always required
	if cfg.Feature == "" {
		color.Red("Error: -f <feature> is required.\n")
		fmt.Println("Run 'netforce --help' for usage.")
		os.Exit(1)
	}

	// Validate the feature name early
	if err := validate.Feature(cfg.Feature); err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}

	// Domain is required for all features except explain and benchmark
	feat := strings.ToLower(cfg.Feature)
	if feat != "explain" && feat != "benchmark" && cfg.Domain == "" {
		color.Red("Error: -d <domain> is required for feature %q.\n", cfg.Feature)
		fmt.Println("Example: netforce -d example.com -f stress -r 50 --duration 30")
		os.Exit(1)
	}

	core.Dispatch(cfg)
}

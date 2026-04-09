package analyzer

import "github.com/lucasenlucas/netforce/internal/metrics"

// DetectRateLimit checks if any response returned HTTP 429 Too Many Requests
func DetectRateLimit(results []metrics.Result) bool {
	for _, r := range results {
		if r.StatusCode == 429 {
			return true
		}
	}
	return false
}

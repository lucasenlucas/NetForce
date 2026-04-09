package core

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/lucasenlucas/netforce/internal/metrics"
)

// worker sends HTTP GET requests in a loop until the context is cancelled
func worker(ctx context.Context, client *http.Client, url string, pacer Pacer, col *metrics.Collector) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Wait for the pacer to permit the next request
			if err := pacer.Wait(ctx); err != nil {
				return
			}

			result := doRequest(ctx, client, url)
			col.Add(result)
		}
	}
}

func doRequest(ctx context.Context, client *http.Client, url string) metrics.Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return metrics.Result{StatusCode: 0, Duration: time.Since(start), Error: err.Error(), Timestamp: start}
	}

	req.Header.Set("User-Agent", "NetForce/1.0 (Performance Testing Tool)")

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return metrics.Result{StatusCode: 0, Duration: duration, Error: err.Error(), Timestamp: start}
	}

	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) // drain to reuse connection

	return metrics.Result{StatusCode: resp.StatusCode, Duration: duration, Timestamp: start}
}

func buildClient(threads int, timeout int) *http.Client {
	dialer := &net.Dialer{
		Timeout:   time.Duration(timeout) * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   threads * 10,
		MaxConnsPerHost:       threads * 10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
}

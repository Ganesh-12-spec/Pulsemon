package checker

import (
	"context"
	"net/http"
	"time"
)

func Check(url string, timeout time.Duration) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()

	_, err = http.DefaultClient.Do(req)

	latency := time.Since(start)

	if err != nil {
		return latency, err
	}

	return latency, nil
}

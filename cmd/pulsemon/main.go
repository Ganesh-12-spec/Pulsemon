package main

import (
	"fmt"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
	"github.com/Ganesh-12-spec/pulsemon/internal/config"
	"github.com/Ganesh-12-spec/pulsemon/internal/history"
)

func main() {
	targets := []config.Target{
		{
			Name: "Example",
			URL:  "https://example.com",
		},
		{
			Name: "Google",
			URL:  "https://google.com",
		},
	}

	// History stores the result of every health check.
	historyStore := history.History{}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Pulsemon started...")

	for range ticker.C {
		for _, target := range targets {

			// Check the target with a 5-second timeout.
			err := checker.Check(target.URL, 5*time.Second)

			if err != nil {
				fmt.Println("Health check failed for", target.Name, ":", err)

				// Store the failed health check.
				historyStore.Add(history.Record{
					Target: target.Name,
					Status: "DOWN",
					Time:   time.Now(),
				})

				continue
			}

			fmt.Println("Health check passed for", target.Name)

			// Store the successful health check.
			historyStore.Add(history.Record{
				Target: target.Name,
				Status: "UP",
				Time:   time.Now(),
			})
		}
	}
}

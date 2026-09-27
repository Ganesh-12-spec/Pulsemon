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

	// History stores every health-check result.
	historyStore := history.History{}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Pulsemon started...")

	for range ticker.C {
		for _, target := range targets {

			latency, err := checker.Check(target.URL, 5*time.Second)

			if err != nil {
				fmt.Println("Health check failed for", target.Name, ":", err)

				historyStore.Add(history.Record{
					Target:  target.Name,
					Status:  "DOWN",
					Time:    time.Now(),
					Latency: latency,
				})

				continue
			}

			fmt.Println("Health check passed for", target.Name)

			historyStore.Add(history.Record{
				Target:  target.Name,
				Status:  "UP",
				Time:    time.Now(),
				Latency: latency,
			})

			uptime := historyStore.Uptime(target.Name)
			averageLatency := historyStore.AverageLatency(target.Name)

			fmt.Printf(
				"%s uptime: %.2f%% | average latency: %v\n",
				target.Name,
				uptime,
				averageLatency,
			)
		}
	}
}

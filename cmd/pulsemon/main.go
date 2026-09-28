package main

import (
	"fmt"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
	"github.com/Ganesh-12-spec/pulsemon/internal/config"
	"github.com/Ganesh-12-spec/pulsemon/internal/history"
	"github.com/Ganesh-12-spec/pulsemon/internal/state"
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

	historyStore := history.History{}
	stateMonitor := state.NewMonitor()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Pulsemon started...")

	for range ticker.C {
		for _, target := range targets {
			latency, err := checker.Check(target.URL, 5*time.Second)

			status, changed := stateMonitor.Update(target.Name, err)

			if err != nil {
				fmt.Println("Health check failed for", target.Name, ":", err)

				historyStore.Add(history.Record{
					Target:  target.Name,
					Status:  string(status),
					Time:    time.Now(),
					Latency: latency,
				})

				if changed {
					fmt.Println(target.Name, "state changed to", status)
				}

				continue
			}

			fmt.Println("Health check passed for", target.Name)

			historyStore.Add(history.Record{
				Target:  target.Name,
				Status:  string(status),
				Time:    time.Now(),
				Latency: latency,
			})

			if changed {
				fmt.Println(target.Name, "state changed to", status)
			}

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

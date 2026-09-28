package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
	"github.com/Ganesh-12-spec/pulsemon/internal/config"
	"github.com/Ganesh-12-spec/pulsemon/internal/history"
	"github.com/Ganesh-12-spec/pulsemon/internal/state"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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

	logger.Info("Pulsemon started")

	for range ticker.C {
		for _, target := range targets {
			latency, err := checker.Check(target.URL, 5*time.Second)

			status, changed, thresholdReached, recovered := stateMonitor.Update(target.Name, err)

			if err != nil {
				logger.Error(
					"health check failed",
					"target", target.Name,
					"status", status,
					"error", err,
					"latency", latency,
				)

				historyStore.Add(history.Record{
					Target:  target.Name,
					Status:  string(status),
					Time:    time.Now(),
					Latency: latency,
				})

				if changed {
					logger.Warn(
						"target state changed",
						"target", target.Name,
						"status", status,
					)
				}

				if thresholdReached {
					logger.Error(
						"failure threshold reached",
						"target", target.Name,
						"status", status,
						"consecutive_failures", 3,
					)
				}

				continue
			}

			logger.Info(
				"health check passed",
				"target", target.Name,
				"status", status,
				"latency", latency,
			)

			historyStore.Add(history.Record{
				Target:  target.Name,
				Status:  string(status),
				Time:    time.Now(),
				Latency: latency,
			})

			if changed {
				logger.Info(
					"target state changed",
					"target", target.Name,
					"status", status,
				)
			}
			if recovered {
				logger.Info(
					"target recovered",
					"target", target.Name,
					"status", status,
				)
			}

			uptime := historyStore.Uptime(target.Name)
			averageLatency := historyStore.AverageLatency(target.Name)

			logger.Info(
				"target metrics",
				"target", target.Name,
				"uptime", uptime,
				"average_latency", averageLatency,
			)
		}
	}
}

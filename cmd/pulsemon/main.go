package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/api"
	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
	"github.com/Ganesh-12-spec/pulsemon/internal/config"
	"github.com/Ganesh-12-spec/pulsemon/internal/history"
	"github.com/Ganesh-12-spec/pulsemon/internal/metrics"
	"github.com/Ganesh-12-spec/pulsemon/internal/state"
	"github.com/Ganesh-12-spec/pulsemon/internal/webhook"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	webhookURL := "https://example.com/webhook"

	targets := []config.Target{
		{Name: "Example", URL: "https://example.com"},
		{Name: "Google", URL: "https://google.com"},
	}

	historyStore := history.History{}
	stateMonitor := state.NewMonitor()
	metricsStore := metrics.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/status", api.StatusHandler(stateMonitor))
	mux.HandleFunc("/metrics", metricsStore.Handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Create a context that is cancelled when Ctrl+C or SIGTERM is received.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Start the HTTP API in the background.
	go func() {
		logger.Info("status API started", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("status API stopped", "error", err)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	logger.Info("Pulsemon started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutdown signal received")
			goto shutdown

		case <-ticker.C:
			for _, target := range targets {
				latency, err := checker.Check(target.URL, 5*time.Second)

				status, changed, thresholdReached, recovered := stateMonitor.Update(
					target.Name,
					err,
				)

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

						err := webhook.Send(webhookURL, webhook.Payload{
							Target:              target.Name,
							Status:              string(status),
							Event:               "failure_threshold",
							ConsecutiveFailures: 3,
						})

						if err != nil {
							logger.Error(
								"failed to send failure webhook",
								"target", target.Name,
								"error", err,
							)
						}
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

					err := webhook.Send(webhookURL, webhook.Payload{
						Target: target.Name,
						Status: string(status),
						Event:  "recovered",
					})

					if err != nil {
						logger.Error(
							"failed to send recovery webhook",
							"target", target.Name,
							"error", err,
						)
					}
				}

				uptime := historyStore.Uptime(target.Name)
				averageLatency := historyStore.AverageLatency(target.Name)

				metricsStore.Update(
					target.Name,
					string(status),
					uptime,
					averageLatency.Seconds()*1000,
				)

				logger.Info(
					"target metrics",
					"target", target.Name,
					"uptime", uptime,
					"average_latency", averageLatency,
				)
			}
		}
	}

shutdown:

	// Give active HTTP requests up to 5 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		return
	}

	logger.Info("Pulsemon stopped cleanly")
}

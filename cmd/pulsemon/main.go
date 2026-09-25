package main

import (
	"fmt"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
	"github.com/Ganesh-12-spec/pulsemon/internal/config"
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

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Pulsemon started...")

	for range ticker.C {
		for _, target := range targets {
			err := checker.Check(target.URL)

			if err != nil {
				fmt.Println("Health check failed for", target.Name, ":", err)
				continue
			}

			fmt.Println("Health check passed for", target.Name)
		}
	}
}

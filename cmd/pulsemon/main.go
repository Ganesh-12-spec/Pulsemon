package main

import (
	"fmt"

	"github.com/Ganesh-12-spec/pulsemon/internal/config"
)

func main() {
	// Define one sample monitoring target.
	target := config.Target{
		Name:      "Example API",
		URL:       "https://example.com",
		Interval:  10,
		Timeout:   5,
		Threshold: 3,
	}

	// Print the configuration temporarily.
	// Actual health-checking logic will be added
	// in a later commit.
	fmt.Printf("Monitoring target: %+v\n", target)
}

package main

import (
	"fmt"
	"time"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Pulsemon started...")

	for range ticker.C {
		err := checker.Check("https://example.com")

		if err != nil {
			fmt.Println("Health check failed:", err)
			continue
		}

		fmt.Println("Health check passed")
	}
}

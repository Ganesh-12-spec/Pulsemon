package main

import (
	"fmt"

	"github.com/Ganesh-12-spec/pulsemon/internal/checker"
)

func main() {
	err := checker.Check("https://example.com")

	if err != nil {
		fmt.Println("Health check failed:", err)
		return
	}

	fmt.Println("Health check passed")
}

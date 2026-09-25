package config

import "fmt"

type Target struct {
	Name      string
	URL       string
	Interval  int
	Timeout   int
	Threshold int
}

func (t Target) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if t.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	if t.Interval <= 0 {
		return fmt.Errorf("interval must be a positive integer")
	}
	if t.Timeout <= 0 {
		return fmt.Errorf("timeout must be a positive integer")
	}
	if t.Threshold < 0 {
		return fmt.Errorf("threshold cannot be negative")
	}
	return nil
}

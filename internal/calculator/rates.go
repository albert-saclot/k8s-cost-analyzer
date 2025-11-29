package calculator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Rates struct {
	CPUPerCorePerHour    float64 `yaml:"cpu_per_core_per_hour"`
	MemoryPerGBPerHour   float64 `yaml:"memory_per_gb_per_hour"`
}

// DefaultRates returns reasonable default pricing
func DefaultRates() Rates {
	return Rates{
		CPUPerCorePerHour:  0.034,
		MemoryPerGBPerHour: 0.004,
	}
}

// LoadRatesFromFile reads rates from a YAML file
func LoadRatesFromFile(path string) (Rates, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Rates{}, fmt.Errorf("failed to read rates file: %w", err)
	}

	var rates Rates
	if err := yaml.Unmarshal(data, &rates); err != nil {
		return Rates{}, fmt.Errorf("failed to parse rates YAML: %w", err)
	}

	return rates, nil
}

// GetRatesLastUpdated extracts the last updated date from rates file
// Returns the date and number of days since update, or error if not found
func GetRatesLastUpdated(path string) (time.Time, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("failed to open rates file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Last updated:") {
			parts := strings.SplitN(line, "Last updated:", 2)
			if len(parts) != 2 {
				continue
			}
			dateStr := strings.TrimSpace(parts[1])
			lastUpdated, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return time.Time{}, 0, fmt.Errorf("invalid date format: %w", err)
			}
			daysSince := int(time.Since(lastUpdated).Hours() / 24)
			return lastUpdated, daysSince, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, 0, fmt.Errorf("error reading file: %w", err)
	}

	return time.Time{}, 0, fmt.Errorf("no 'Last updated' date found in rates file")
}

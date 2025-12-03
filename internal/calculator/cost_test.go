package calculator

import (
	"math"
	"os"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	// tolerance allows for floating-point precision errors in cost calculations
	tolerance = 0.0001
)

func TestCalculatePodCost(t *testing.T) {
	t.Parallel()
	rates := Rates{
		CPUPerCorePerHour:  0.034,
		MemoryPerGBPerHour: 0.004,
	}

	tests := []struct {
		name           string
		podName        string
		namespace      string
		cpuRequest     string
		memoryRequest  string
		expectedHourly float64
		expectedDaily  float64
		expectedMonth  float64
	}{
		{
			name:           "typical pod with 100m CPU and 128Mi memory",
			podName:        "test-pod",
			namespace:      "default",
			cpuRequest:     "100m",
			memoryRequest:  "128Mi",
			expectedHourly: 0.0034 + (128.0 / 1024 * 0.004), // 0.1 cores * 0.034 + 0.125 GB * 0.004
			expectedDaily:  (0.0034 + (128.0 / 1024 * 0.004)) * 24,
			expectedMonth:  (0.0034 + (128.0 / 1024 * 0.004)) * 730,
		},
		{
			name:           "pod with 1 CPU and 1Gi memory",
			podName:        "big-pod",
			namespace:      "default",
			cpuRequest:     "1",
			memoryRequest:  "1Gi",
			expectedHourly: 0.034 + 0.004, // 1 core * 0.034 + 1 GB * 0.004
			expectedDaily:  (0.034 + 0.004) * 24,
			expectedMonth:  (0.034 + 0.004) * 730,
		},
		{
			name:           "pod with zero requests",
			podName:        "zero-pod",
			namespace:      "default",
			cpuRequest:     "0",
			memoryRequest:  "0",
			expectedHourly: 0,
			expectedDaily:  0,
			expectedMonth:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cpuQty := resource.MustParse(tt.cpuRequest)
			memQty := resource.MustParse(tt.memoryRequest)

			cost := CalculatePodCost(tt.podName, tt.namespace, cpuQty, memQty, rates)

			if cost.Name != tt.podName {
				t.Errorf("expected name %s, got %s", tt.podName, cost.Name)
			}

			if cost.Namespace != tt.namespace {
				t.Errorf("expected namespace %s, got %s", tt.namespace, cost.Namespace)
			}

			if math.Abs(cost.Hourly.TotalCost-tt.expectedHourly) > tolerance {
				t.Errorf("hourly cost: expected %.4f, got %.4f", tt.expectedHourly, cost.Hourly.TotalCost)
			}

			if math.Abs(cost.Daily.TotalCost-tt.expectedDaily) > tolerance {
				t.Errorf("daily cost: expected %.2f, got %.2f", tt.expectedDaily, cost.Daily.TotalCost)
			}

			if math.Abs(cost.Monthly.TotalCost-tt.expectedMonth) > tolerance {
				t.Errorf("monthly cost: expected %.2f, got %.2f", tt.expectedMonth, cost.Monthly.TotalCost)
			}

			// Verify breakdown consistency
			expectedHourlyCPU := float64(cpuQty.MilliValue()) / 1000.0 * rates.CPUPerCorePerHour
			if math.Abs(cost.Hourly.CPUCost-expectedHourlyCPU) > tolerance {
				t.Errorf("hourly CPU cost: expected %.4f, got %.4f", expectedHourlyCPU, cost.Hourly.CPUCost)
			}

			expectedHourlyMemory := float64(memQty.Value()) / (1024 * 1024 * 1024) * rates.MemoryPerGBPerHour
			if math.Abs(cost.Hourly.MemoryCost-expectedHourlyMemory) > tolerance {
				t.Errorf("hourly memory cost: expected %.4f, got %.4f", expectedHourlyMemory, cost.Hourly.MemoryCost)
			}
		})
	}
}

func TestDefaultRates(t *testing.T) {
	t.Parallel()
	rates := DefaultRates()

	if rates.CPUPerCorePerHour <= 0 {
		t.Errorf("expected positive CPU rate, got %.4f", rates.CPUPerCorePerHour)
	}

	if rates.MemoryPerGBPerHour <= 0 {
		t.Errorf("expected positive memory rate, got %.4f", rates.MemoryPerGBPerHour)
	}
}

func TestLoadRatesFromFile(t *testing.T) {
	t.Parallel()
	t.Run("valid file", func(t *testing.T) {
		t.Parallel()
		tmpfile, err := os.CreateTemp("", "rates-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		content := `cpu_per_core_per_hour: 0.05
memory_per_gb_per_hour: 0.006
`
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpfile.Close()

		rates, err := LoadRatesFromFile(tmpfile.Name())
		if err != nil {
			t.Fatalf("LoadRatesFromFile failed: %v", err)
		}

		if rates.CPUPerCorePerHour != 0.05 {
			t.Errorf("expected CPU rate 0.05, got %.4f", rates.CPUPerCorePerHour)
		}

		if rates.MemoryPerGBPerHour != 0.006 {
			t.Errorf("expected memory rate 0.006, got %.4f", rates.MemoryPerGBPerHour)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		t.Parallel()
		_, err := LoadRatesFromFile("/nonexistent/file.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		t.Parallel()
		tmpfile, err := os.CreateTemp("", "rates-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte("invalid: yaml: content:\n  bad")); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpfile.Close()

		_, err = LoadRatesFromFile(tmpfile.Name())
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})
}

func TestGetRatesLastUpdated(t *testing.T) {
	t.Parallel()
	t.Run("valid date", func(t *testing.T) {
		t.Parallel()
		tmpfile, err := os.CreateTemp("", "rates-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		testDate := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
		content := `# Default cost rates (USD)
# Last updated: 2023-01-15
#
cpu_per_core_per_hour: 0.034
memory_per_gb_per_hour: 0.004
`
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpfile.Close()

		lastUpdated, daysSince, err := GetRatesLastUpdated(tmpfile.Name())
		if err != nil {
			t.Fatalf("GetRatesLastUpdated failed: %v", err)
		}

		if lastUpdated.Year() != 2023 || lastUpdated.Month() != time.January || lastUpdated.Day() != 15 {
			t.Errorf("expected date 2023-01-15, got %v", lastUpdated)
		}

		expectedDaysSince := int(time.Since(testDate).Hours() / 24)
		diff := daysSince - expectedDaysSince
		if diff < -1 || diff > 1 {
			t.Errorf("expected daysSince around %d (±1), got %d", expectedDaysSince, daysSince)
		}
	})

	t.Run("no date found", func(t *testing.T) {
		t.Parallel()
		tmpfile, err := os.CreateTemp("", "rates-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		content := `cpu_per_core_per_hour: 0.034
memory_per_gb_per_hour: 0.004
`
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpfile.Close()

		_, _, err = GetRatesLastUpdated(tmpfile.Name())
		if err == nil {
			t.Error("expected error when no date found, got nil")
		}
	})

	t.Run("invalid file", func(t *testing.T) {
		t.Parallel()
		_, _, err := GetRatesLastUpdated("/nonexistent/file.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})
}

package reporter

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

// Note: This test cannot use t.Parallel() because captureStdout modifies os.Stdout,
// which is global state. Running these tests in parallel would cause interference.
func TestPrintCostCSV(t *testing.T) {
	tests := []struct {
		name              string
		costs             []calculator.PodCost
		wantLineCount     int
		wantFirstPodName  string
		wantContainsCosts []string
		wantSecondPodName string
	}{
		{
			name: "single pod",
			costs: []calculator.PodCost{
				{
					Name: "test-pod-1",
					Hourly: calculator.ResourceCost{
						CPUCost:    0.01,
						MemoryCost: 0.005,
						TotalCost:  0.015,
					},
					Daily: calculator.ResourceCost{
						CPUCost:    0.24,
						MemoryCost: 0.12,
						TotalCost:  0.36,
					},
					Monthly: calculator.ResourceCost{
						CPUCost:    7.30,
						MemoryCost: 3.65,
						TotalCost:  10.95,
					},
				},
			},
			wantLineCount:     2,
			wantFirstPodName:  "test-pod-1",
			wantContainsCosts: []string{"0.0100", "10.95"},
		},
		{
			name:          "empty input",
			costs:         []calculator.PodCost{},
			wantLineCount: 1,
		},
		{
			name: "multiple pods",
			costs: []calculator.PodCost{
				{
					Name:    "pod-1",
					Hourly:  calculator.ResourceCost{TotalCost: 0.01},
					Daily:   calculator.ResourceCost{TotalCost: 0.24},
					Monthly: calculator.ResourceCost{TotalCost: 7.30},
				},
				{
					Name:    "pod-2",
					Hourly:  calculator.ResourceCost{TotalCost: 0.02},
					Daily:   calculator.ResourceCost{TotalCost: 0.48},
					Monthly: calculator.ResourceCost{TotalCost: 14.60},
				},
			},
			wantLineCount:     3,
			wantFirstPodName:  "pod-1",
			wantSecondPodName: "pod-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				if err := PrintCostCSV(tt.costs); err != nil {
					t.Fatalf("PrintCostCSV failed: %v", err)
				}
			})

			lines := strings.Split(strings.TrimSpace(output), "\n")

			if len(lines) != tt.wantLineCount {
				t.Fatalf("line count: got %d, want %d", len(lines), tt.wantLineCount)
			}

			// Verify header
			expectedHeader := "pod_name,hourly_cpu_cost,hourly_memory_cost,hourly_total_cost,daily_cpu_cost,daily_memory_cost,daily_total_cost,monthly_cpu_cost,monthly_memory_cost,monthly_total_cost"
			if lines[0] != expectedHeader {
				t.Errorf("header mismatch:\nexpected: %s\ngot:      %s", expectedHeader, lines[0])
			}

			// Parse CSV for proper verification
			if tt.wantLineCount > 1 {
				reader := csv.NewReader(strings.NewReader(output))
				records, err := reader.ReadAll()
				if err != nil {
					t.Fatalf("failed to parse CSV output: %v", err)
				}

				if len(records) != tt.wantLineCount {
					t.Fatalf("CSV record count: got %d, want %d", len(records), tt.wantLineCount)
				}

				// Verify first data row
				if tt.wantFirstPodName != "" {
					firstRow := records[1]
					if firstRow[0] != tt.wantFirstPodName {
						t.Errorf("first pod name: got %q, want %q", firstRow[0], tt.wantFirstPodName)
					}
				}

				// Verify cost values in output
				for _, costValue := range tt.wantContainsCosts {
					if !strings.Contains(output, costValue) {
						t.Errorf("expected output to contain cost value %q", costValue)
					}
				}

				// Verify second data row if present
				if tt.wantSecondPodName != "" {
					secondRow := records[2]
					if secondRow[0] != tt.wantSecondPodName {
						t.Errorf("second pod name: got %q, want %q", secondRow[0], tt.wantSecondPodName)
					}
				}
			}
		})
	}
}

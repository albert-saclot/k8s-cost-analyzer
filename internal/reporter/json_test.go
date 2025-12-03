package reporter

import (
	"encoding/json"
	"testing"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

// Note: This test cannot use t.Parallel() because captureStdout modifies os.Stdout,
// which is global state. Running these tests in parallel would cause interference.
func TestPrintCostJSON(t *testing.T) {
	tests := []struct {
		name              string
		namespace         string
		costs             []calculator.PodCost
		wantNamespace     string
		wantPodCount      int
		wantFirstPodName  string
		wantHourlyTotal   float64
		wantMonthlyCost   float64
		wantSummaryPods   int
		wantSummaryMonthly float64
	}{
		{
			name:      "single pod",
			namespace: "default",
			costs: []calculator.PodCost{
				{
					Name:      "test-pod-1",
					Namespace: "default",
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
						CPUCost:    7.3,
						MemoryCost: 3.65,
						TotalCost:  10.95,
					},
				},
			},
			wantNamespace:      "default",
			wantPodCount:       1,
			wantFirstPodName:   "test-pod-1",
			wantHourlyTotal:    0.015,
			wantMonthlyCost:    10.95,
			wantSummaryPods:    1,
			wantSummaryMonthly: 10.95,
		},
		{
			name:               "empty input",
			namespace:          "default",
			costs:              []calculator.PodCost{},
			wantNamespace:      "default",
			wantPodCount:       0,
			wantSummaryPods:    0,
			wantSummaryMonthly: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				if err := PrintCostJSON(tt.namespace, tt.costs); err != nil {
					t.Fatalf("PrintCostJSON failed: %v", err)
				}
			})

			// Parse JSON output
			var jsonOut jsonOutput
			if err := json.Unmarshal([]byte(output), &jsonOut); err != nil {
				t.Fatalf("failed to parse JSON output: %v", err)
			}

			// Verify structure
			if jsonOut.Namespace != tt.wantNamespace {
				t.Errorf("namespace: got %q, want %q", jsonOut.Namespace, tt.wantNamespace)
			}

			if len(jsonOut.Pods) != tt.wantPodCount {
				t.Fatalf("pod count: got %d, want %d", len(jsonOut.Pods), tt.wantPodCount)
			}

			if tt.wantPodCount > 0 {
				pod := jsonOut.Pods[0]
				if pod.Name != tt.wantFirstPodName {
					t.Errorf("first pod name: got %q, want %q", pod.Name, tt.wantFirstPodName)
				}

				if pod.Hourly.TotalCost != tt.wantHourlyTotal {
					t.Errorf("hourly total: got %.4f, want %.4f", pod.Hourly.TotalCost, tt.wantHourlyTotal)
				}
			}

			// Verify summary
			if jsonOut.Summary.TotalPods != tt.wantSummaryPods {
				t.Errorf("summary total_pods: got %d, want %d", jsonOut.Summary.TotalPods, tt.wantSummaryPods)
			}

			if jsonOut.Summary.MonthlyCost != tt.wantSummaryMonthly {
				t.Errorf("summary monthly_cost: got %.2f, want %.2f", jsonOut.Summary.MonthlyCost, tt.wantSummaryMonthly)
			}
		})
	}
}

package analyzer

import (
	"math"
	"testing"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

const (
	// hourlyTolerance allows for floating-point precision errors in hourly calculations
	hourlyTolerance = 0.0001

	// dailyMonthlyTolerance is larger due to accumulation of floating-point errors
	// when multiplying hourly rates by 24 or 730
	dailyMonthlyTolerance = 0.01
)

func TestAggregateByNamespace(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		costs         []calculator.PodCost
		wantNamespace string
		wantPods      int
		wantHourly    float64
		wantDaily     float64
		wantMonthly   float64
	}{
		{
			name: "multiple pods",
			costs: []calculator.PodCost{
				{
					Name:      "pod1",
					Namespace: "default",
					Hourly:    calculator.ResourceCost{CPUCost: 0.01, MemoryCost: 0.005, TotalCost: 0.015},
					Daily:     calculator.ResourceCost{CPUCost: 0.24, MemoryCost: 0.12, TotalCost: 0.36},
					Monthly:   calculator.ResourceCost{CPUCost: 7.3, MemoryCost: 3.65, TotalCost: 10.95},
				},
				{
					Name:      "pod2",
					Namespace: "default",
					Hourly:    calculator.ResourceCost{CPUCost: 0.02, MemoryCost: 0.01, TotalCost: 0.03},
					Daily:     calculator.ResourceCost{CPUCost: 0.48, MemoryCost: 0.24, TotalCost: 0.72},
					Monthly:   calculator.ResourceCost{CPUCost: 14.6, MemoryCost: 7.3, TotalCost: 21.9},
				},
			},
			wantNamespace: "default",
			wantPods:      2,
			wantHourly:    0.045,
			wantDaily:     1.08,
			wantMonthly:   32.85,
		},
		{
			name:          "empty input",
			costs:         []calculator.PodCost{},
			wantNamespace: "",
			wantPods:      0,
			wantHourly:    0,
			wantDaily:     0,
			wantMonthly:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			summary := AggregateByNamespace(tt.costs)

			if summary.Namespace != tt.wantNamespace {
				t.Errorf("namespace: got %q, want %q", summary.Namespace, tt.wantNamespace)
			}

			if summary.TotalPods != tt.wantPods {
				t.Errorf("total pods: got %d, want %d", summary.TotalPods, tt.wantPods)
			}

			if math.Abs(summary.HourlyCost-tt.wantHourly) > hourlyTolerance {
				t.Errorf("hourly cost: got %.4f, want %.4f", summary.HourlyCost, tt.wantHourly)
			}

			if math.Abs(summary.DailyCost-tt.wantDaily) > dailyMonthlyTolerance {
				t.Errorf("daily cost: got %.2f, want %.2f", summary.DailyCost, tt.wantDaily)
			}

			if math.Abs(summary.MonthlyCost-tt.wantMonthly) > dailyMonthlyTolerance {
				t.Errorf("monthly cost: got %.2f, want %.2f", summary.MonthlyCost, tt.wantMonthly)
			}
		})
	}
}

func TestSortByMonthlyCost(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		costs     []calculator.PodCost
		wantOrder []string // expected pod names in order
	}{
		{
			name: "three pods",
			costs: []calculator.PodCost{
				{Name: "cheap-pod", Monthly: calculator.ResourceCost{TotalCost: 5.0}},
				{Name: "expensive-pod", Monthly: calculator.ResourceCost{TotalCost: 20.0}},
				{Name: "medium-pod", Monthly: calculator.ResourceCost{TotalCost: 10.0}},
			},
			wantOrder: []string{"expensive-pod", "medium-pod", "cheap-pod"},
		},
		{
			name:      "empty input",
			costs:     []calculator.PodCost{},
			wantOrder: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			originalFirst := ""
			if len(tt.costs) > 0 {
				originalFirst = tt.costs[0].Name
			}

			sorted := SortByMonthlyCost(tt.costs)

			if len(sorted) != len(tt.wantOrder) {
				t.Fatalf("length: got %d, want %d", len(sorted), len(tt.wantOrder))
			}

			for i, wantName := range tt.wantOrder {
				if sorted[i].Name != wantName {
					t.Errorf("position %d: got %q, want %q", i, sorted[i].Name, wantName)
				}
			}

			// Verify original slice not modified
			if len(tt.costs) > 0 && tt.costs[0].Name != originalFirst {
				t.Errorf("original slice was modified: first element is now %q, was %q", tt.costs[0].Name, originalFirst)
			}
		})
	}
}

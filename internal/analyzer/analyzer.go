package analyzer

import (
	"sort"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

type NamespaceSummary struct {
	Namespace    string
	TotalPods    int
	HourlyCost   float64
	DailyCost    float64
	MonthlyCost  float64
}

// AggregateByNamespace sums up costs for all pods in a namespace
func AggregateByNamespace(costs []calculator.PodCost) NamespaceSummary {
	if len(costs) == 0 {
		return NamespaceSummary{}
	}

	summary := NamespaceSummary{
		Namespace: costs[0].Namespace,
		TotalPods: len(costs),
	}

	for _, pc := range costs {
		summary.HourlyCost += pc.Hourly.TotalCost
		summary.DailyCost += pc.Daily.TotalCost
		summary.MonthlyCost += pc.Monthly.TotalCost
	}

	return summary
}

// SortByMonthlyCost sorts pod costs by monthly total (descending)
func SortByMonthlyCost(costs []calculator.PodCost) []calculator.PodCost {
	sorted := make([]calculator.PodCost, len(costs))
	copy(sorted, costs)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Monthly.TotalCost > sorted[j].Monthly.TotalCost
	})

	return sorted
}

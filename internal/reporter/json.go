package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

type jsonOutput struct {
	Namespace string              `json:"namespace"`
	Pods      []jsonPodCost       `json:"pods"`
	Summary   jsonNamespaceSummary `json:"summary"`
}

type jsonPodCost struct {
	Name    string           `json:"name"`
	Hourly  jsonResourceCost `json:"hourly"`
	Daily   jsonResourceCost `json:"daily"`
	Monthly jsonResourceCost `json:"monthly"`
}

type jsonResourceCost struct {
	CPUCost    float64 `json:"cpu_cost"`
	MemoryCost float64 `json:"memory_cost"`
	TotalCost  float64 `json:"total_cost"`
}

type jsonNamespaceSummary struct {
	TotalPods   int     `json:"total_pods"`
	HourlyCost  float64 `json:"hourly_cost"`
	DailyCost   float64 `json:"daily_cost"`
	MonthlyCost float64 `json:"monthly_cost"`
}

// PrintCostJSON outputs pod costs in JSON format
func PrintCostJSON(namespace string, costs []calculator.PodCost) error {
	pods := make([]jsonPodCost, len(costs))
	var totalHourly, totalDaily, totalMonthly float64

	for i, c := range costs {
		pods[i] = jsonPodCost{
			Name: c.Name,
			Hourly: jsonResourceCost{
				CPUCost:    c.Hourly.CPUCost,
				MemoryCost: c.Hourly.MemoryCost,
				TotalCost:  c.Hourly.TotalCost,
			},
			Daily: jsonResourceCost{
				CPUCost:    c.Daily.CPUCost,
				MemoryCost: c.Daily.MemoryCost,
				TotalCost:  c.Daily.TotalCost,
			},
			Monthly: jsonResourceCost{
				CPUCost:    c.Monthly.CPUCost,
				MemoryCost: c.Monthly.MemoryCost,
				TotalCost:  c.Monthly.TotalCost,
			},
		}
		totalHourly += c.Hourly.TotalCost
		totalDaily += c.Daily.TotalCost
		totalMonthly += c.Monthly.TotalCost
	}

	output := jsonOutput{
		Namespace: namespace,
		Pods:      pods,
		Summary: jsonNamespaceSummary{
			TotalPods:   len(costs),
			HourlyCost:  totalHourly,
			DailyCost:   totalDaily,
			MonthlyCost: totalMonthly,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

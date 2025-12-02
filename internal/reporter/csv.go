package reporter

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

// PrintCostCSV outputs pod costs in CSV format
func PrintCostCSV(costs []calculator.PodCost) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// Write header
	header := []string{
		"pod_name",
		"hourly_cpu_cost",
		"hourly_memory_cost",
		"hourly_total_cost",
		"daily_cpu_cost",
		"daily_memory_cost",
		"daily_total_cost",
		"monthly_cpu_cost",
		"monthly_memory_cost",
		"monthly_total_cost",
	}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, c := range costs {
		row := []string{
			c.Name,
			fmt.Sprintf("%.4f", c.Hourly.CPUCost),
			fmt.Sprintf("%.4f", c.Hourly.MemoryCost),
			fmt.Sprintf("%.4f", c.Hourly.TotalCost),
			fmt.Sprintf("%.2f", c.Daily.CPUCost),
			fmt.Sprintf("%.2f", c.Daily.MemoryCost),
			fmt.Sprintf("%.2f", c.Daily.TotalCost),
			fmt.Sprintf("%.2f", c.Monthly.CPUCost),
			fmt.Sprintf("%.2f", c.Monthly.MemoryCost),
			fmt.Sprintf("%.2f", c.Monthly.TotalCost),
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

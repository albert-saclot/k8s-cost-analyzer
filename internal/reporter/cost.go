package reporter

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
)

// PrintCostTable displays pod costs in a formatted table
func PrintCostTable(costs []calculator.PodCost) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "POD\tHOURLY\tDAILY\tMONTHLY")
	for _, c := range costs {
		fmt.Fprintf(w, "%s\t$%.4f\t$%.2f\t$%.2f\n",
			c.Name,
			c.Hourly.TotalCost,
			c.Daily.TotalCost,
			c.Monthly.TotalCost,
		)
	}
}

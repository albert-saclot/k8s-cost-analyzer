package reporter

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/k8s"
)

// PrintPodResourcesTable displays pod resources in a formatted table
func PrintPodResourcesTable(resources []k8s.PodResources) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "POD\tCPU REQUEST\tMEMORY REQUEST\tCPU LIMIT\tMEMORY LIMIT")
	for _, r := range resources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			r.Name, r.CPURequest, r.MemoryRequest, r.CPULimit, r.MemoryLimit)
	}
}

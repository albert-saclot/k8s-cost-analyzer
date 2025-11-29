package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/analyzer"
	"github.com/albert-saclot/k8s-cost-analyzer/internal/calculator"
	"github.com/albert-saclot/k8s-cost-analyzer/internal/k8s"
	"github.com/albert-saclot/k8s-cost-analyzer/internal/reporter"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze resource requests and costs for a namespace",
	Long:  `Display pod resource requests, limits, and estimated costs for the specified namespace.`,
	RunE:  runAnalyze,
}

var (
	namespace  string
	cpuRate    float64
	memoryRate float64
	showCosts  bool
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to analyze")
	analyzeCmd.Flags().Float64Var(&cpuRate, "cpu-rate", 0.034, "Cost per CPU core per hour (USD)")
	analyzeCmd.Flags().Float64Var(&memoryRate, "memory-rate", 0.004, "Cost per GB memory per hour (USD)")
	analyzeCmd.Flags().BoolVar(&showCosts, "costs", true, "Show cost estimates")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// Check if using default rates and if they're stale
	usingDefaultCPU := !cmd.Flags().Changed("cpu-rate")
	usingDefaultMemory := !cmd.Flags().Changed("memory-rate")
	if showCosts && (usingDefaultCPU || usingDefaultMemory) {
		checkRateStaleness()
	}

	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ctx := context.Background()
	pods, err := k8s.FetchPods(ctx, client, namespace)
	if err != nil {
		return err
	}

	if len(pods) == 0 {
		fmt.Printf("No pods found in namespace '%s'\n", namespace)
		return nil
	}

	fmt.Printf("Analyzing %d pods in namespace '%s':\n\n", len(pods), namespace)

	if !showCosts {
		resources := make([]k8s.PodResources, len(pods))
		for i, pod := range pods {
			resources[i] = k8s.ExtractResources(pod)
		}
		reporter.PrintPodResourcesTable(resources)
		return nil
	}

	// Calculate costs
	rates := calculator.Rates{
		CPUPerCorePerHour:  cpuRate,
		MemoryPerGBPerHour: memoryRate,
	}

	podCosts := make([]calculator.PodCost, 0, len(pods))
	for _, pod := range pods {
		res := k8s.ExtractResources(pod)

		cpuQty, _ := resource.ParseQuantity(res.CPURequest)
		memQty, _ := resource.ParseQuantity(res.MemoryRequest)

		if cpuQty.IsZero() && memQty.IsZero() {
			continue
		}

		cost := calculator.CalculatePodCost(pod.Name, pod.Namespace, cpuQty, memQty, rates)
		podCosts = append(podCosts, cost)
	}

	// Sort by cost (highest first)
	sortedCosts := analyzer.SortByMonthlyCost(podCosts)

	reporter.PrintCostTable(sortedCosts)

	// Show summary
	summary := analyzer.AggregateByNamespace(podCosts)
	fmt.Printf("\nNamespace Summary:\n")
	fmt.Printf("  Total Pods: %d\n", summary.TotalPods)
	fmt.Printf("  Estimated Monthly Cost: $%.2f\n", summary.MonthlyCost)
	fmt.Printf("\nNote: These are estimates based on resource requests, not actual usage.\n")

	return nil
}

func checkRateStaleness() {
	ratesPath := "config/rates.yaml"
	_, daysSince, err := calculator.GetRatesLastUpdated(ratesPath)
	if err != nil {
		// Silently ignore if rates file doesn't exist or can't be read
		return
	}

	const staleThreshold = 180 // 6 months
	if daysSince > staleThreshold {
		fmt.Fprintf(os.Stderr, "Warning: Default pricing rates are %d days old (last updated %d days ago).\n", daysSince, daysSince)
		fmt.Fprintf(os.Stderr, "Consider updating %s or use --cpu-rate and --memory-rate flags.\n\n", ratesPath)
	}
}

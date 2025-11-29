package cmd

import (
	"context"
	"fmt"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/k8s"
	"github.com/albert-saclot/k8s-cost-analyzer/internal/reporter"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze resource requests and costs for a namespace",
	Long:  `Display pod resource requests and limits for the specified namespace.`,
	RunE:  runAnalyze,
}

var namespace string

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to analyze")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
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

	resources := make([]k8s.PodResources, len(pods))
	for i, pod := range pods {
		resources[i] = k8s.ExtractResources(pod)
	}

	reporter.PrintPodResourcesTable(resources)
	return nil
}

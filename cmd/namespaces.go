package cmd

import (
	"context"
	"fmt"

	"github.com/albert-saclot/k8s-cost-analyzer/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List all namespaces in the cluster",
	Long:  `Connect to the Kubernetes cluster and list all available namespaces.`,
	RunE:  runNamespaces,
}

func init() {
	rootCmd.AddCommand(namespacesCmd)
}

func runNamespaces(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	ctx := context.Background()
	namespaces, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	fmt.Printf("Found %d namespaces:\n\n", len(namespaces.Items))
	for _, ns := range namespaces.Items {
		fmt.Printf("  %s (Status: %s)\n", ns.Name, ns.Status.Phase)
	}

	return nil
}

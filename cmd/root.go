package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kcost",
	Short: "Analyze Kubernetes resource costs",
	Long: `A CLI tool to analyze and report on Kubernetes resource costs.

This tool connects to your Kubernetes cluster and calculates the cost
of running your workloads based on CPU and memory requests.`,
}

// Execute runs the root command and handles errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "kcost: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// TODO: global flags (--kubeconfig, --context)
}

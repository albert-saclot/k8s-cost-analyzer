package calculator

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceCost struct {
	CPUCost    float64
	MemoryCost float64
	TotalCost  float64
}

type PodCost struct {
	Name      string
	Namespace string
	Hourly    ResourceCost
	Daily     ResourceCost
	Monthly   ResourceCost
}

// CalculatePodCost computes cost for a pod's resource requests
func CalculatePodCost(podName, namespace string, cpuRequest, memoryRequest resource.Quantity, rates Rates) PodCost {
	// Convert CPU to cores (millicores to cores)
	cpuCores := float64(cpuRequest.MilliValue()) / 1000.0

	// Convert memory to GB
	memoryGB := float64(memoryRequest.Value()) / (1024 * 1024 * 1024)

	// Calculate hourly costs
	hourlyCPU := cpuCores * rates.CPUPerCorePerHour
	hourlyMemory := memoryGB * rates.MemoryPerGBPerHour
	hourlyTotal := hourlyCPU + hourlyMemory

	// Calculate daily costs (24 hours)
	dailyCPU := hourlyCPU * 24
	dailyMemory := hourlyMemory * 24
	dailyTotal := hourlyTotal * 24

	// Calculate monthly costs (730 hours average)
	monthlyCPU := hourlyCPU * 730
	monthlyMemory := hourlyMemory * 730
	monthlyTotal := hourlyTotal * 730

	return PodCost{
		Name:      podName,
		Namespace: namespace,
		Hourly: ResourceCost{
			CPUCost:    hourlyCPU,
			MemoryCost: hourlyMemory,
			TotalCost:  hourlyTotal,
		},
		Daily: ResourceCost{
			CPUCost:    dailyCPU,
			MemoryCost: dailyMemory,
			TotalCost:  dailyTotal,
		},
		Monthly: ResourceCost{
			CPUCost:    monthlyCPU,
			MemoryCost: monthlyMemory,
			TotalCost:  monthlyTotal,
		},
	}
}

# k8s-cost-analyzer

A zero-infrastructure CLI tool for quick Kubernetes cost estimates.

Unlike kubectl-cost (which requires OpenCost + Prometheus backend) or full monitoring platforms like Kubecost, kcost is a standalone binary that calculates cost estimates directly from resource requests. Perfect for pre-deployment analysis, on-prem environments, and multi-cluster cost checks without installing monitoring infrastructure.

## Installation

### From source

```bash
git clone https://github.com/albert-saclot/k8s-cost-analyzer.git
cd k8s-cost-analyzer
go build -o kcost
```

### Using go install

```bash
go install github.com/albert-saclot/k8s-cost-analyzer@latest
```

## Usage

### List namespaces

```bash
kcost namespaces
```

### Analyze namespace costs

```bash
kcost analyze -n kube-system
```

Example output:
```
Analyzing 7 pods in namespace 'kube-system':

POD                                HOURLY    DAILY   MONTHLY
kube-apiserver-minikube            $0.0085   $0.20   $6.21
kube-controller-manager-minikube   $0.0068   $0.16   $4.96
etcd-minikube                      $0.0038   $0.09   $2.77

Namespace Summary:
  Total Pods: 5
  Estimated Monthly Cost: $19.10

Note: These are estimates based on resource requests, not actual usage.
```

### Custom pricing rates

```bash
kcost analyze -n production --cpu-rate 0.05 --memory-rate 0.006
```

### View resources without costs

```bash
kcost analyze -n kube-system --costs=false
```

### Output formats

The tool supports multiple output formats for integration with other tools and workflows.

**JSON output:**
```bash
kcost analyze -n kube-system -o json
```

Example output:
```json
{
  "namespace": "kube-system",
  "pods": [
    {
      "name": "kube-apiserver-minikube",
      "hourly": {
        "cpu_cost": 0.0068,
        "memory_cost": 0.0017,
        "total_cost": 0.0085
      },
      "daily": {
        "cpu_cost": 0.16,
        "memory_cost": 0.04,
        "total_cost": 0.20
      },
      "monthly": {
        "cpu_cost": 4.96,
        "memory_cost": 1.25,
        "total_cost": 6.21
      }
    }
  ],
  "summary": {
    "total_pods": 5,
    "hourly_cost": 0.026,
    "daily_cost": 0.63,
    "monthly_cost": 19.10
  }
}
```

**CSV output:**
```bash
kcost analyze -n kube-system -o csv > costs.csv
```

Outputs a CSV file with columns: `pod_name`, hourly/daily/monthly costs for CPU, memory, and total.

### Updating pricing rates

The tool uses default rates from `config/rates.yaml` based on AWS m5.large pricing. These rates are date-stamped and should be reviewed periodically.

**Check current rates:**
```bash
cat config/rates.yaml
```

**Update rates manually:**
1. Edit `config/rates.yaml`
2. Update `cpu_per_core_per_hour` and `memory_per_gb_per_hour` values
3. Update the `Last updated:` date

**Or use the update helper:**
```bash
./scripts/update-rates.sh  # Shows pricing sources and calculation guide
```

The tool warns when default rates are more than 6 months old. Override with `--cpu-rate` and `--memory-rate` flags to suppress warnings.

The tool connects to your Kubernetes cluster using `~/.kube/config`.

## Cost Estimates Disclaimer

**kcost provides cost estimates for planning purposes only.**

These estimates are:
- Based on resource requests, not actual usage or cloud bills
- Calculated using configurable rates that you set
- Excluding storage, networking, load balancers, and other non-compute costs
- Not accounting for reserved instances, spot pricing, or volume discounts

**Good for:**
- Quick cost estimates without installing monitoring infrastructure
- Pre-deployment cost modeling (analyze YAML before applying)
- On-prem clusters with known hardware costs
- Multi-cluster cost checks
- CI/CD cost gates

**Not suitable for:**
- Actual usage tracking (use kubectl-cost + OpenCost for this)
- Historical analysis and trends
- Production billing or financial reporting

Cost estimates may differ significantly from actual cloud bills. No warranty or liability is provided for financial decisions based on this tool's output.

For production cost monitoring with usage tracking, see [OpenCost](https://github.com/opencost/opencost) or [kubectl-cost](https://github.com/kubecost/kubectl-cost).

## Requirements

- Go 1.23+
- Access to a Kubernetes cluster
- Valid kubeconfig file

## Development

### Testing locally

If you don't have a cluster, use minikube:

```bash
brew install minikube
minikube start
./kcost namespaces
```

### Building

```bash
go build -o kcost
```

### Project Structure

```
k8s-cost-analyzer/
├── main.go                  # Entry point
├── cmd/                     # Cobra commands
│   ├── root.go             # Root command
│   ├── namespaces.go       # Namespace listing
│   └── analyze.go          # Cost analysis
├── internal/
│   ├── k8s/                # Kubernetes client
│   ├── calculator/         # Cost calculation
│   ├── analyzer/           # Cost aggregation
│   └── reporter/           # Output formatting
├── config/
│   └── rates.yaml          # Default pricing rates
└── scripts/
    └── update-rates.sh     # Pricing update helper
```

## Roadmap

- [x] Basic CLI structure and K8s connection
- [x] Fetch and display pod resource requests/limits
- [x] Cost calculation engine with configurable rates
- [x] Multiple output formats (JSON, CSV)
- [ ] Testing and documentation
- [ ] Multi-namespace analysis
- [ ] Resource usage analysis (via metrics-server)
- [ ] Cost optimization recommendations

## License

MIT

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

List all namespaces in your cluster:

```bash
kcost namespaces
```

Example output:
```
Found 4 namespaces:

  default (Status: Active)
  kube-node-lease (Status: Active)
  kube-public (Status: Active)
  kube-system (Status: Active)
```

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
├── main.go              # Entry point
├── cmd/                 # Cobra commands
│   ├── root.go         # Root command
│   └── namespaces.go   # Namespace listing
└── internal/
    └── k8s/
        └── client.go   # K8s client wrapper
```

## Roadmap

- [x] Basic CLI structure and K8s connection
- [ ] Fetch and display pod resource requests/limits
- [ ] Cost calculation engine with configurable rates
- [ ] Multiple output formats (JSON, CSV, table)
- [ ] Resource usage analysis
- [ ] Cost optimization recommendations

## License

MIT

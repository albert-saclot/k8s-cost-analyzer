#!/bin/bash
# Update pricing rates in config/rates.yaml
#
# CURRENT STATUS: Manual process documented below
# FUTURE: Could be automated to fetch from cloud provider APIs

set -e

RATES_FILE="config/rates.yaml"
TODAY=$(date +%Y-%m-%d)

echo "Pricing Rate Update Tool"
echo "========================"
echo ""
echo "Current process (manual):"
echo "1. Check cloud provider pricing:"
echo "   - AWS: https://aws.amazon.com/ec2/pricing/on-demand/"
echo "   - GCP: https://cloud.google.com/compute/vm-instance-pricing"
echo "   - Azure: https://azure.microsoft.com/en-us/pricing/details/virtual-machines/"
echo ""
echo "2. Select a reference instance (e.g., m5.large on AWS)"
echo "3. Calculate per-core and per-GB rates:"
echo "   - CPU rate = instance_price / num_cores"
echo "   - Memory rate = instance_price / memory_gb"
echo "   - Apply blending for typical request/limit patterns"
echo ""
echo "4. Update $RATES_FILE:"
echo "   - Set cpu_per_core_per_hour"
echo "   - Set memory_per_gb_per_hour"
echo "   - Update 'Last updated' date to $TODAY"
echo ""
echo "Future enhancement ideas:"
echo "- Fetch AWS pricing via AWS Price List API"
echo "- Query GCP Cloud Billing Catalog API"
echo "- Pull from OpenCost's maintained pricing data"
echo "- Support multiple regions/instance types"
echo ""
echo "For now, please update $RATES_FILE manually."

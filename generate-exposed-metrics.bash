#!/bin/bash
set -eo pipefail
oldIFS=$IFS
IFS=$'\n'
echo "Metric name | Description"
echo "---- | ----------"
for line in $(grep "metric: Create" src/metrics/*); do
  metric_name=$(echo "$line" | cut -d\" -f2)
  metric_desc=$(echo "$line" | cut -d\" -f4)
  echo "$metric_name | $metric_desc"
done
IFS=$oldIFS


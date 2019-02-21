#!/bin/bash

set -m

echo never | tee /sys/kernel/mm/transparent_hugepage/enabled
echo never | tee /sys/kernel/mm/transparent_hugepage/defrag

# Start the first process
redis-server &

# Start the second process
./goack
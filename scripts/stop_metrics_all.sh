#!/bin/bash

NUM_SERVERS=${1:-11}

for i in $(seq 2 $NUM_SERVERS); do
  echo "Stopping send_metrics.sh in fake-server-$i ..."
  docker exec fake-server-$i pkill -f send_metrics.sh || true
done

echo "All send_metrics.sh processes stopped!"

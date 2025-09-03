#!/bin/bash

NUM_SERVERS=${1:-11}

for i in $(seq 2 $NUM_SERVERS); do
  docker exec -d fake-server-$i bash send_metrics.sh &
done

wait
echo "All send_metrics processes started in parallel!"

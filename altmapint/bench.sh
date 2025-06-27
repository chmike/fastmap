#!/bin/bash

# verify that an argument is provided
if [ $# -eq 0 ]; then
    echo "error: $0 requires the benchmark selector as argument"
    echo "usage: $0 <BENCH_PART>"
    echo "example: $0 Cache2Hit"
    exit 1
fi

RUNS=10
BENCH_NAME=$(basename "$(pwd)")
BENCH_PART="$1"
GOARCH=$(go env GOARCH)
FNAME="stats_$GOARCH.txt"

rm -f "$FNAME"

for i in $(seq 1 $RUNS); do
	echo "$BENCH_NAME:$BENCH_PART $i"
    go test -bench="$BENCH_PART" | sed "s/$BENCH_NAME/map/g" >> "$FNAME"
done


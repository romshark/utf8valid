#!/bin/sh

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <FILTER> <COUNT>"
  echo "Example 1: $0 . 10"
  echo "Example 2: $0 ValidString 10"
  exit 1
fi

FILTER=$1
COUNT=$2

# Replace any '/' with '_'.
FILTER_AS_FILENAME="${FILTER//\//_}"

OUT_STD=".std_$FILTER_AS_FILENAME.txt"
OUT_OPTIMIZED=".opt_$FILTER_AS_FILENAME.txt"
FUNC_STD="std"
FUNC_OPTIMIZED="optimized"

echo "Benchmarking function $FUNC_STD to $OUT_STD"
go test \
  -test.timeout=0 \
  -benchmem \
  -bench $FILTER \
  -benchfunc $FUNC_STD \
  -count $COUNT \
  | tee $OUT_STD
clear


echo "Benchmarking function $FUNC_OPTIMIZED to $OUT_OPTIMIZED"
go test \
  -test.timeout=0 \
  -benchmem \
  -bench $FILTER \
  -benchfunc $FUNC_OPTIMIZED \
  -count $COUNT \
  | tee $OUT_OPTIMIZED
clear

go run golang.org/x/perf/cmd/benchstat@latest $OUT_STD $OUT_OPTIMIZED
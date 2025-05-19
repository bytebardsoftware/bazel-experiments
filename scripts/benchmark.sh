#!/bin/bash
set -euo pipefail
set -x

bazel run //:bazel-experiments_test --  -test.count 1 -test.run=^# -test.bench=. -test.benchtime=5s
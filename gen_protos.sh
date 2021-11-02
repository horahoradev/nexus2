#!/bin/bash
set -euo pipefail

docker build -t grpcutil .
docker run -v $(pwd)/multiplayerservice:/multiplayerservice -it -t grpcutil
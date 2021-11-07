#!/bin/bash
set -euo pipefail

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative ./multiplayerservice/protocol/multiplayerservice.proto
python3 -m grpc_tools.protoc -I ./multiplayerservice/protocol/ --python_out=./multiplayerservice/protocol/ --grpc_python_out=./multiplayerservice/protocol/ multiplayerservice.proto
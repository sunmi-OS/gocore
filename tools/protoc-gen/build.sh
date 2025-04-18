#!/bin/bash
#
# lxy@20230907
#
set -ve
cd cmd/protoc-gen-sm-openapi/ && go build && cd ../../
cd cmd/protoc-gen-sm-go-errors/ && go build && cd ../../
cd cmd/protoc-gen-sm-go-gin/ && go build && cd ../../


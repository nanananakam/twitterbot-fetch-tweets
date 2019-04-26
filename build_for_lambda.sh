#!/usr/bin/env bash

# macだとうまくbuildできないのでdocker内とか
GOOS=linux GOARCH=amd64 go build -o Main
zip Main.zip Main
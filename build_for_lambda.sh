#!/usr/bin/env bash

docker run -e GO111MODULE=on -v $(pwd):/go/src/github.com/nanananakam/twitterbot-fetch-tweets golang:1.12-stretch go build -o /go/src/github.com/nanananakam/twitterbot-fetch-tweets/Main /go/src/github.com/nanananakam/twitterbot-fetch-tweets/main.go
zip Main.zip Main
rm go.sum
rm -f Main
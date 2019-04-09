#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 go build -o Main
zip Main.zip Main
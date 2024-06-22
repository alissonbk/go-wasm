#!/bin/zsh

kill -9 "$(ps -ef | grep "python -m http.server" | grep -v grep | awk '{print $2}')"; GOOS=js GOARCH=wasm go build -o main.wasm ./main.go; python -m http.server

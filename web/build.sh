#!/bin/sh
cd $(dirname $0)
GOOS=js GOARCH=wasm go build -o ./Raylib-Go-Wasm/index/main.wasm .

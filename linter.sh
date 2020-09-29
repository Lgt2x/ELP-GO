#!/bin/bash

gofmt -w src/elp/main.go
goimports -w src/elp/main.go
#!/bin/bash
git checkout *
git pull
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main ../jx3/
cd ../jx3/
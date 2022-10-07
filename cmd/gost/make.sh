#!/bin/bash

date --rfc-3339=seconds
echo "" #隔行

rm gost-*

# linux x86 32bit
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o gost-linux-x86
# linux x86 64bit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gost-linux-x64
# 适应华硕AC86U路由器
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o gost-linux-arm
# 适应腾达AC18路由器
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o gost-linux-armv5

# windows x86 32bit
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gost-win-x86.exe
# windows x86 64bit
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gost-win-x64.exe

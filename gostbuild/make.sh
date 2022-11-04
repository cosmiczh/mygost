#!/bin/bash

date --rfc-3339=seconds
echo "" #隔行

rm gost-*

# gost-win-x86.exe
#echo gost-win-x86.exe
#CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gost-win-x86.exe
# windows x86 64bit
echo gost-win-x64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gost-win-x64.exe ../cmd/gost

# gost-linux-x86
#echo gost-linux-x86
#CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o gost-linux-x86
# gost-linux-x64
echo gost-linux-x64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gost-linux-x64 ../cmd/gost

# gost-linux-arm64
echo gost-linux-arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gost-linux-arm64 ../cmd/gost
# gost-linux-armv5 适应腾达AC18路由器
echo gost-linux-armv5
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o gost-linux-armv5 ../cmd/gost
# gost-linux-armv6
echo gost-linux-armv6
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o gost-linux-armv6 ../cmd/gost
# gost-linux-armv7 适应华硕AC86U路由器
echo gost-linux-armv7
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o gost-linux-armv7 ../cmd/gost

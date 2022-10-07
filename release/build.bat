set GOOS=windows
set GOARCH=amd64
go build -o gost-win-x64.exe ../cmd/gost

set GOOS=linux
go build -o gost-linux-x64 ../cmd/gost
set GOARCH=arm
go build -o gost-linux-arm ../cmd/gost


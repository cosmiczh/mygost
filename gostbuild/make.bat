del gost-*

:: gost-win-x86.exe
::echo gost-win-x86.exe
::set CGO_ENABLED=0
::set GOOS=windows
::set GOARCH=386
::go build -o gost-win-x86.exe  ../cmd/gost

:: gost-win-x64.exe
echo gost-win-x64.exe
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -o gost-win-x64.exe  ../cmd/gost

:: gost-linux-x86
::echo gost-linux-x86
::set CGO_ENABLED=0
::set GOOS=linux
::set GOARCH=386
::go build -o gost-linux-x86  ../cmd/gost

:: gost-linux-x64
echo gost-linux-x64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o gost-linux-x64  ../cmd/gost

:: gost-linux-arm64
echo gost-linux-arm64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -o gost-linux-arm64  ../cmd/gost
:: gost-linux-armv5 适应腾达AC18路由器
echo gost-linux-armv5
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm
set GOARM=5
go build -o gost-linux-armv5  ../cmd/gost
:: gost-linux-armv6
echo gost-linux-armv6
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm
set GOARM=6
go build -o gost-linux-armv6  ../cmd/gost
:: gost-linux-armv7 适应华硕AC86U路由器
echo gost-linux-armv7
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o gost-linux-armv7  ../cmd/gost

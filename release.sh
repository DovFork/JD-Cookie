#!/bin/bash
DRONE_TAG=v3.0.1
mkdir -f /tmp/build/
rsync -av --exclude build --exclude .github --exclude .idea  ./ /tmp/build/
cd /tmp/build
go get -u github.com/gobuffalo/packr/v2/packr2
go mod tidy
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=amd64   packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_amd64
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=386           packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_i386
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=arm   GOARM=7 packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo  -o dist/jdcookie_v3_linux_armv7
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=arm64         packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo  -o dist/jdcookie_v3_linux_arm64
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=ppc64         packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_ppc64
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=ppc64le       packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_ppc64le
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=mips          packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_mips
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=mipsle        packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_mipsle
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=mips64        packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_mips64
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=linux GOARCH=mips64le      packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_linux_mips64le
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=windows GOARCH=386         packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo  -o dist/jdcookie_v3_windows_i386.exe
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=windows GOARCH=amd64       packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_windows_adm64.exe
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=windows GOARCH=arm GOARM=7 packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_windows_arm.exe
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=darwin GOARCH=arm64        packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_darwin_arm64
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=darwin GOARCH=amd64        packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_darwin_amd64
#CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=android GOARCH=arm   GOARM=7      packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_android_arm
CGO_LDFLAGS="-static" CGO_ENABLE=1 GOOS=android GOARCH=arm64       packr2 build  -ldflags   "-w -s -X main.Build=`date +%FT%T%z` -X main.Version=$DRONE_TAG" -installsuffix cgo -o dist/jdcookie_v3_android_arm64

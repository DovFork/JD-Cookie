[![Build Status](https://drone.pi.scjtqs.com:8443/api/badges/scjtqs/jd_cookie/status.svg)](https://drone.pi.scjtqs.com:8443/scjtqs/jd_cookie)
# 说明

## v1.x 版本(tag) 本地运行，本地提取
京东的cookie 本地提取工具(小白专用) 

双击exe文件，运行服务，后用本地的浏览器打开 http://127.0.0.1:29099

来辅助提取你的cookie

## v2.x 版本(tag) 服务器运行，客户端自行扫码提取 
> 已不受服务端和客户端在同ip的限制
> 
> 客户端通过浏览器打开， http://公网ip或者域名:29099/
> 
> 然后扫码即可提取。
>
> 
```bash
## linux/macos下
chmod +x jdcookie
## 加上-d 代表后台运行
./jdcookie -d
## windows 下
## 可以直接双击运行。也可以cmd中运行
## cmd 中运行 参数和linux一样
.\jdcookie.exe -d
```


## 编译方法
```bash
#安装编译工具，用于打包html到二进制
go get -u github.com/gobuffalo/packr/v2/packr2
#build和go build一样，交叉编译的方式也是一样的
CGO_ENABLED=0 GOOS=linux GOARCH=arm packr2 build  -o jdcookie
#简单点的
packr2 build  -o jdcookie
# 编译windows x86版本
CGO_ENABLED=0 GOOS=windows GOARCH=386 packr2 build  -o jdcookie_windows_x86.exe
# 编译windows x64版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 packr2 build  -o jdcookie_windows_x64.exe
# 编译mac arm64 (m1)版本
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 packr2 build  -o jdcookie_darwin_arm64
# 编译mac amd64 x64版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 packr2 build  -o jdcookie_darwin_x64
```
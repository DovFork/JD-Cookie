BUILD=`date +%FT%T%z`
VER=v3.0.1

LDFLAGS=-ldflags " -s -X main.Build=${BUILD} -X main.Version=${Version}"

build :
	rm -rf dist
	mkdir dist
	CGO_ENABLED=1 CGO_LDFLAGS="-static" packr2 build  ${LDFLAGS} -installsuffix cgo -o ./dist/jdcookie .
	chmod -R +x ./dist

clean:
	rm -rf dist

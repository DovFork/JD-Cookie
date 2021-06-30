BUILD=`date +%FT%T%z`
VER=v2.0.1

LDFLAGS=-ldflags " -s -X main.Build=${BUILD} -X main.Version=${Version}"

build :
	rm -rf dist
	mkdir dist
	packr2 build  ${LDFLAGS} -o ./dist/jdcookie .
	chmod -R +x ./dist

clean:
	rm -rf dist

BUILD=`date +%FT%T%z`
VER=v2.0.1

LDFLAGS=-a -ldflags " -s -X main.Build=${BUILD} -X main.Version=${VER} -X main.Version=${Version} -extldflags '-static'"

build :
	rm -rf dist
	mkdir dist
	packr2 build  ${LDFLAGS} -o ./dist/jdcookie .
	chmod -R +x ./dist

clean:
	rm -rf dist

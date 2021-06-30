FROM golang:1.16-alpine as builder

LABEL name="jd_cookie server"
LABEL version="2.0.1"
LABEL author="scjtqs <scjtqs@qq.com>"
LABEL maintainer="scjtqs <scjtqs@qq.com>"
LABEL description="simple to get jd cookie"

ARG Version="v2.0.1"

ADD . /src
ENV GOPROXY "http://goproxy.cn,direct"
ENV CGO_ENABLED "0"
ENV GO111MODULE "on"

ENV UPSAVE ""

ENV DB_ENABLE "false"
ENV DB_HOST ""
ENV DB_PORT ""
ENV DB_USER ""
ENV DB_PASS ""
ENV DB_DATABASE ""

RUN cd /src \
    && apk add --no-cache  make \
    && rm -rf dist \
    && go get -u github.com/gobuffalo/packr/v2/packr2 \
    && go mod tidy \
    && make

FROM alpine:3.13 as production

COPY --from=builder /src/dist /opt/app

RUN  adduser -D -H www \
     && chown -R www /opt/app

USER www
WORKDIR /opt/app

EXPOSE 29099

CMD ["/opt/app/jdcookie"]
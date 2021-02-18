#docker build --network host --rm --build-arg APP_ROOT=/go/src/kubecp -t kubecp:latest -f Dockerfile .
#0 ----------------------------
FROM golang:1.15.3
ARG  APP_ROOT
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}

ENV GO111MODULE=on
#ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
ENV PATH=$GOPATH/bin:$PATH

# install upx
RUN sed -i "s/deb.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list \
  && sed -i "s/security.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list \
  && apt-get update \
  && apt-get install upx musl-dev -y

# build code
RUN go get -u github.com/swaggo/swag/cmd/swag \
  && swag init \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o main . \
  && strip --strip-unneeded main \
  && upx --lzma main

#1 ----------------------------
FROM alpine:latest
ARG APP_ROOT
WORKDIR /app
COPY --from=0 ${APP_ROOT}/savior .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk add --no-cache openssh jq curl busybox-extras \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/app/main"]

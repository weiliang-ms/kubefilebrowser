#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
#go get -u github.com/swaggo/swag/cmd/swag
swag init
# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
  echo "!!!!!!Swagger documentation generate error, please check the source code!!!!!!"
  exit 1
fi

# build web
cd web && yarn run build && cd ../
sed -i "s/Vue App/KubeFileBrowser/g" static/index.html
# build server
name="kubefilebrowser"
version="v1.6.1"
# linux
archList="386 amd64 arm arm64 ppc64le"
# shellcheck disable=SC2181
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_linux-"$i"-"$version
  CGO_ENABLED=0 GOOS=linux GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done

# windows
# shellcheck disable=SC2181
archList="386 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=BinaryName=$name"_windows-"$i"-"$version".exe"
  CGO_ENABLED=0 GOOS=windows GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done

# darwin
# shellcheck disable=SC2181
archList="arm64 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=BinaryName=$name"_darwin-"$i"-"$version
  CGO_ENABLED=0 GOOS=darwin GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done
#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
#go get -u github.com/swaggo/swag/cmd/swag
swag init
if [ "$?" != "0" ]; then
  echo "!!!!!!Swagger documentation generate error, please check the source code!!!!!!"
  exit 1
fi
name="kubecp"
version="1.0-bate"
osList="linux windows darwin"
# shellcheck disable=SC2181
for i in $osList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_"$i"-"$version
  if [ "$i" == "windows" ];then
    BinaryName=${BinaryName}.exe
  fi
  CGO_ENABLED=0 GOOS=$i GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o "$BinaryName" .
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Compilation error, please check the source code!!!!!!"
    exit 1
  fi

#  strip --strip-unneeded $BinaryName
#  # shellcheck disable=SC2181
#  if [ "$?" != "0" ]; then
#    echo "!!!!!!Strip failed!!!!!!"
#    exit 1
#  fi

  upx --lzma "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Upx failed!!!!!!"
    exit 1
  fi
done

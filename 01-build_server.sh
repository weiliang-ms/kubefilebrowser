#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
#go get -u github.com/swaggo/swag/cmd/swag
swag init -g cmd/server/main.go
# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
  echo "!!!!!!Swagger documentation generate error, please check the source code!!!!!!"
  exit 1
fi

# build web
cd web && yarn run build && cd ../

# build server
name="kubefilebrowser"
version="v1.3-beta"
osList="linux windows darwin"
# shellcheck disable=SC2181
for i in $osList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_"$i"-"$version
  if [ "$i" == "windows" ];then
    BinaryName=${BinaryName}.exe
  fi
  CGO_ENABLED=0 GOOS=$i GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Server compilation error, please check the source code!!!!!!"
    exit 1
  fi

#  strip --strip-unneeded "$BinaryName"
#  # shellcheck disable=SC2181
#  if [ "$?" != "0" ]; then
#    echo "!!!!!!Strip server failed!!!!!!"
#    exit 1
#  fi

  upx --lzma "$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Upx server failed!!!!!!"
    exit 1
  fi
done

# shellcheck disable=SC2154
# shellcheck disable=SC2027
tar zcf $name"_"$version".tar.gz" static $name"_linux-"$version $name"_windows-"$version".exe" $name"_darwin-"$version
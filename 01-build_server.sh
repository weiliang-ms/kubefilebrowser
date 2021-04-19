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
version="v1.6"
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
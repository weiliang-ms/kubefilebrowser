#!/usr/bin/env bash

mkdir -p utils/ls_binary
osList="linux windows"
# shellcheck disable=SC2181
for i in $osList; do
  # shellcheck disable=SC2027
  BinaryName="ls_"$i"_amd64"
  if [ "$i" == "windows" ];then
    # shellcheck disable=SC2027
    BinaryName="ls_"$i"_amd64".exe
  fi
  CGO_ENABLED=0 GOOS=$i GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o utils/ls_binary/"$BinaryName" cmd/ls/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/ls_binary/"$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Upx ls failed!!!!!!"
    exit 1
  fi
done

mkdir -p utils/zip_binary
osList="linux windows"
# shellcheck disable=SC2181
for i in $osList; do
  # shellcheck disable=SC2027
  BinaryName="zip_"$i"_amd64"
  if [ "$i" == "windows" ];then
    # shellcheck disable=SC2027
    BinaryName="zip_"$i"_amd64".exe
  fi
  CGO_ENABLED=0 GOOS=$i GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o utils/zip_binary/"$BinaryName" cmd/zip/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi

  upx --lzma utils/zip_binary/"$BinaryName"
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!Upx ls failed!!!!!!"
    exit 1
  fi
done

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
#sed -i "s/Vue App/KubeFileBrowser/g" static/index.html
find static -name "index.html" -type f | xargs sed -i".bak" "s/Vue App/KubeFileBrowser/g" && find static -name "*.bak" -exec rm -rf {} \;
# build server
CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o kubefilebrowser
# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
  echo "!!!!!!Server compilation error, please check the source code!!!!!!"
  exit 1
fi
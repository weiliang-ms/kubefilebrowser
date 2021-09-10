#!/usr/bin/env bash

mkdir -p utils/kf_tools_binary
# linux
archList="386 amd64 arm arm64 ppc64le"
# shellcheck disable=SC2181
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kf_tools_linux_"$i
  CGO_ENABLED=0 GOOS=linux GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kf_tools_binary/"$BinaryName"
done

# windows
# shellcheck disable=SC2181
archList="386 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kf_tools_windows_"$i".exe"
  CGO_ENABLED=0 GOOS=windows GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kf_tools_binary/"$BinaryName"
done

# darwin
# shellcheck disable=SC2181
archList="arm64 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kf_tools_darwin_"$i
  CGO_ENABLED=0 GOOS=darwin GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kf_tools_binary/"$BinaryName"
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
#cd web && yarn run build && cd ../
#sed -i "s/Vue App/KubeFileBrowser/g" static/index.html
# shellcheck disable=SC2038
find static -name "index.html" -type f | xargs sed -i".bak" "s/Vue App/KubeFileBrowser/g" && find static -name "*.bak" -exec rm -rf {} \;
# build server
CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o kubefilebrowser
# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
  echo "!!!!!!Server compilation error, please check the source code!!!!!!"
  exit 1
fi
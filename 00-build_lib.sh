#!/usr/bin/env bash

mkdir -p utils/kf_tools_binary
# linux
archList="386 amd64 arm arm64 ppc64le"
# shellcheck disable=SC2181
BuildAt=$(date)
GitHash=$(git rev-parse --short HEAD)
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kf_tools_linux_"$i
  CGO_ENABLED=0 GOOS=linux GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
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
  CGO_ENABLED=0 GOOS=windows GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
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
  CGO_ENABLED=0 GOOS=darwin GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kf_tools_binary/"$BinaryName" cmd/kf_tools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kf_tools_binary/"$BinaryName"
done
#!/bin/sh


build_binary(){
    local binary_name="$1"
    echo building "$binary_name" ...
    go get -v
    go build -v -ldflags "-X 'main.version=$VERSION' -X 'main.commitHash=$SHA1' -X 'main.buildDate=$TIMESTAMP' \
-X 'main.goOS=$GOOS' -X 'main.goArch=$GOARCH'" -o "../$binary_name"
}


[ -z "$GOPATH" ] && export GOPATH=$HOME/go

export VERSION="`git describe --tags`"
export SHA1="`git rev-parse HEAD`"
#export TIMESTAMP="`date -u -Iseconds`"
export TIMESTAMP="`date -u '+%FT%T_%Z'`"

#export PKG_NAME='gitlab.rayark.com/fieliapm/distributeur_automatique'
#export PKG_PATH="$GOPATH/src/$PKG_NAME"

export GO111MODULE='on'
export GOPRIVATE='*.rayark.com'
# turn off proxy and sumdb only if we are doing testing
#export GOPROXY='direct'
#export GOSUMDB='off'

export CGO_ENABLED='0'

# change to project cmd directory
old_dir="`pwd`"
cd cmd

export GOOS='windows'
export GOARCH='amd64'
build_binary 'distributeur_automatique-windows-amd64.exe'
unset GOOS
unset GOARCH

export GOOS='darwin'
export GOARCH='amd64'
build_binary 'distributeur_automatique-macosx-amd64'
unset GOOS
unset GOARCH

export GOOS='linux'
export GOARCH='amd64'
build_binary 'distributeur_automatique-linux-amd64'
unset GOOS
unset GOARCH

# change dir to project root directory
cd "$old_dir"


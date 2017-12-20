#!/bin/bash

set -o pipefail
set -o errexit

PROJECT_DIR=$(pwd)

generate_version() {
        local gitVersion=$(git describe)

        local programVersion=$(cat << VR | sed "s/@@PLACEHOLDER@@/${gitVersion}/g"
package main

var (
        VERSION string = "@@PLACEHOLDER@@"
)
VR
)
        echo "${programVersion}" | tee cmd/client/version.go cmd/server/version.go > /dev/null
}

resolve_deps() {
        go get "golang.org/x/crypto/ssh/terminal"
}

build_go() {
        local path="$1"
        local buildName="$2"

        cd "$path"
        go build -o "${PROJECT_DIR}/build/${buildName}"
        cd "${PROJECT_DIR}"
}

build() {
        if [[ ! -d "build" ]]; then
                mkdir build
        else
                rm -rf build/irc-client
        fi

        resolve_deps
        generate_version

        build_go "cmd/client" "client"
        build_go "cmd/server" "server"

        if [[ ! -e ${PROJECT_DIR}/misc/UserConfigs.json ]]; then
                cp ${PROJECT_DIR}/misc/UserConfigs.json ${PROJECT_DIR}/build
        fi
}

build
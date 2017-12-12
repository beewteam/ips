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
        echo "${programVersion}" > cmd/version.go 
}

resolve_deps() {
        go get "golang.org/x/crypto/ssh/terminal"
}

build() {
        if [[ ! -d "build" ]]; then
                mkdir build
        else
                rm -rf build/irc-client
        fi

        resolve_deps
        generate_version

        cd "cmd"
        go build -o ${PROJECT_DIR}/build/irc-client
        if [[ ! -e ${PROJECT_DIR}/misc/UserConfigs.json ]]; then
                cp ${PROJECT_DIR}/misc/UserConfigs.json ${PROJECT_DIR}/build
        fi
        cd ${PROJECT_DIR}
}

build
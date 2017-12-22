#!/bin/bash

set -o pipefail
set -o errexit

EXTERNAL_DEPENDENCIES=(
        "golang.org/x/crypto/ssh/terminal"
        "github.com/stretchr/testify/assert"
)

COMPONENTS=(
        "cmd/client"
        "cmd/server"
)

PROJECT_DIR=$(pwd)
PROJECT_BUILD_OUTPUT_DIR="build"

VERSION_TEMPLATE=$(cat << VR
package main

var (
        VERSION string = "@@VERSION_PLACEHOLDER@@"
)
VR
)

generate_version() {
        local gitVersion=$(git describe)
        local programVersion=$(echo $VERSION_TEMPLATE | sed "s/@@VERSION_PLACEHOLDER@@/${gitVersion}/g")
        echo "${programVersion}" | tee cmd/client/version.go cmd/server/version.go > /dev/null
}

resolve_deps() {
        for component in ${EXTERNAL_DEPENDENCIES[@]}
        do
                go get $component
        done
}

run_test() {
        for component in ${COMPONENTS[@]}
        do
                cd "${component}"
                go test || exit 1
                cd ${PROJECT_DIR}
        done
}

run_build() {
        if [[ ! -d "build" ]]; then
                mkdir build
        fi

        resolve_deps
        generate_version

        for component in ${COMPONENTS[@]}
        do
                cd "$component"
                go build -o "${PROJECT_DIR}/${PROJECT_BUILD_OUTPUT_DIR}/$(basename ${component})"
                cd "${PROJECT_DIR}"
        done

        #if [[ ! -e ${PROJECT_DIR}/misc/UserConfigs.json ]]; then
        #        cp ${PROJECT_DIR}/misc/UserConfigs.json ${PROJECT_DIR}/build
        #fi
}

run_clean() {
        rm -rf "${PROJECT_BUILD_OUTPUT_DIR}"
}

main() {
        case $1 in
        "-c")
                run_clean
                ;;
        "-t")
                run_test
                ;;
        "-b")
                run_build
                ;;
        *)
                run_clean
                run_build
                run_test
                ;;
        esac
}

main $@
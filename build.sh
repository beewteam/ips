PROJECT_DIR=$(pwd)

build() {
        if [[ ! -e build ]]; then
                mkdir build
        else
                rm -rf build/irc-client
        fi

        cd cmd
        go build -o ${PROJECT_DIR}/build/irc-client
        if [[ ! -e ${PROJECT_DIR}/misc/UserConfigs.json ]]; then
                cp ${PROJECT_DIR}/misc/UserConfigs.json ${PROJECT_DIR}/build
        fi
}

build
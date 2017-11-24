PROJECT_DIR=$(pwd)

build() {
        if [[ ! -e build ]]; then
                mkdir build
        else
                rm -rf build/*
        fi

        cd cmd
        go build -o ${PROJECT_DIR}/build/irc-client
        cp ${PROJECT_DIR}/misc/UserConfigs.json ${PROJECT_DIR}/build
}

build
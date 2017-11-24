build() {
        if [[ ! -e build ]]; then
                mkdir build
        else
                rm -rf build/*
        fi

        cd cmd
        go build -o ../build/irc-client
}

build
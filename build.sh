if [[ ! -e build ]]; then
	mkdir build
fi
cd src
go build -o ../build/irc

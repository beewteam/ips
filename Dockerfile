FROM golang:1.8

WORKDIR /go/src/github.com/beewteam/ips
COPY . .

RUN go get github.com/beewteam/ips; exit 0
RUN ./build.sh
CMD ["build/client", ""]
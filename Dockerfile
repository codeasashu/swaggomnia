FROM golang:1.15
MAINTAINER Fyb3roptik <nwallace@fyberstudios.com>
WORKDIR /go/src/github.com/fyb3roptik/swaggomnia/
COPY . .
RUN go mod init
RUN GO111MODULE=on go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o swaggymnia .
ENTRYPOINT ["./swaggomnia"]

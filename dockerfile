FROM golang:1.15 AS builder
COPY . /go/src/dummy-net-init
WORKDIR /go/src/dummy-net-init
RUN go get &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/src/dummy-net-init/dummy-net-init

FROM scratch
COPY --from='builder' /go/src/dummy-net-init/dummy-net-init /dummy-net-init
ENTRYPOINT ["/dummy-net-init"]
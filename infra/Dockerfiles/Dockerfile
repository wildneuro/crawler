FROM golang:latest as builder

WORKDIR $GOPATH/src/YOURREPO/crawler

COPY . .
RUN go build -mod vendor -o /go/bin/crawler.amd64.bin ./cmd/crawler/crawler.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir -p /voc/bin
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder /go/bin/crawler.amd64.bin /voc/bin/

ENTRYPOINT ["/voc/bin/crawler.amd64.bin"]

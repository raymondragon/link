FROM golang:alpine as builder
WORKDIR /root
ADD . .
WORKDIR /root/src
RUN go mod init link && go mod tidy
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags '-w -s'
FROM scratch
WORKDIR /
COPY --from=builder /root/src/link .
ENTRYPOINT ["/link"]

FROM golang:alpine as builder
WORKDIR /root
ADD . .
RUN go mod init link && go mod tidy
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags '-w -s'
FROM scratch
WORKDIR /
COPY --from=builder /root/link .
ENTRYPOINT ["/link"]
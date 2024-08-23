FROM golang:alpine AS builder
WORKDIR /root
ADD . .
WORKDIR /root/cmd/link
RUN go mod init link && go mod tidy
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags '-w -s'
FROM alpine
COPY --from=builder /root/cmd/link/link /link
ENTRYPOINT ["/link"]

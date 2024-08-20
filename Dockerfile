FROM golang:alpine as builder
WORKDIR /app
ADD . .
RUN go mod init link && go mod tidy && go build -v -trimpath -ldflags '-w -s'
FROM scratch
COPY --from=builder /app/link /link
ENTRYPOINT ["/link"]

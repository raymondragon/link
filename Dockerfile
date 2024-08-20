FROM golang:alpine as builder
WORKDIR /app
ADD . .
RUN cd ./src
RUN go mod init link && go mod tidy
RUN env CGO_ENABLED=0 go build -v -trimpath -ldflags '-w -s'
FROM scratch
COPY --from=builder /app/src/link /link
ENTRYPOINT ["/link"]

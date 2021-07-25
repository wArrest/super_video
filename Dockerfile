FROM golang:1.15.6

# 必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=release
WORKDIR /build
COPY . .
RUN go build -o app .

WORKDIR /dist
RUN cp /build/app .
EXPOSE 8080
CMD ["/dist/app"]

FROM golang:1.23.2-alpine3.20 AS builder

# Define the project name | 定义项目名称

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags '-s -w' \
    -gcflags="all=-trimpath=${PWD}" \
    -asmflags="all=-trimpath=${PWD}" \
    -ldflags="-s -w" \
    -o alertmanager-webhook-adapter ./cmd/alertmanager-webhook-adapter

#linux/amd64,linux/arm64
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/alertmanager-webhook-adapter /bin/alertmanager-webhook-adapter
RUN chmod +x /bin/alertmanager-webhook-adapter && \
    apk update && apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone


CMD ["alertmanager-webhook-adapter","--debug"]

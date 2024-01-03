FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && ls \
    && go mod tidy \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o /build/alertmanager-webhook-adapter cmd/alertmanager-webhook-adapter/main.go

FROM --platform=linux/amd64 alpine:latest

ARG AUTHOR=1228022817@qq.com
LABEL org.opencontainers.image.authors=${AUTHOR}

WORKDIR /app
ENV TZ=Asia/Shanghai
RUN apk update --no-cache && apk add --no-cache tzdata

COPY --from=builder /build/alertmanager-webhook-adapter /usr/local/bin/


EXPOSE 9100

CMD ["alertmanager-webhook-adapter"]

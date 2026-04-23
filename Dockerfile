FROM registry.cn-hangzhou.aliyuncs.com/kubeop/golang:1.25.9 AS builder

COPY . /src/password-self-service

WORKDIR /src/password-self-service

RUN go install github.com/swaggo/swag/cmd/swag@latest && make init && make build

FROM registry.cn-hangzhou.aliyuncs.com/kubeop/alpine:3.23
LABEL maintainer="Sonic.ma<songlin.ma@outlook.com>"

COPY --from=builder /src/password-self-service/password-self-service /usr/local/bin/

RUN apk upgrade --update

USER   nobody
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/password-self-service"]
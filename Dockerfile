FROM docker.mirrors.sjtug.sjtu.edu.cn/kubeop/golang:1.21.4 as BUILD

COPY . /src/password-self-service

WORKDIR /src/password-self-service

RUN go install github.com/swaggo/swag/cmd/swag@latest && make init && make build

FROM docker.mirrors.sjtug.sjtu.edu.cn/kubeop/alpine:3.18
LABEL maintainer="Sonic.ma<songlin.ma@outlook.com>"

COPY --from=BUILD /src/password-self-service/password-self-service /usr/local/bin/

RUN apk upgrade --update

USER   nobody
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/password-self-service"]
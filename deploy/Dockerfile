FROM golang:1.21.4-alpine AS builder

RUN apk add --no-cache bash

LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on
WORKDIR /build
ADD go.mod .
COPY . .

ADD go.sum .
COPY . .

#RUN go mod download
COPY . .

COPY app/interior/api/etc /cmd/etc

RUN  go build -ldflags="-s -w" -o /cmd/protocol app/interior/api/apis.go

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ARG TZ=Asia/Shanghai

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

RUN apk update \
        && apk upgrade \
        && apk add --no-cache bash \
        bash-doc \
        bash-completion \
        && rm -rf /var/cache/apk/* \
        && /bin/bash

WORKDIR /cmd

COPY --from=builder /cmd/protocol /cmd/protocol
COPY --from=builder /cmd/etc /cmd/etc

CMD ["./protocol", "-f", "etc/pro/apis.yaml"]

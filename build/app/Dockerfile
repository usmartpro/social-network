# Собираем в гошке
FROM golang:1.24.1 AS build

ENV BIN_FILE /opt/social/app
ENV CODE_DIR /go/src/
ENV CGO_ENABLED=0
#export GOSUMDB=off
#GOPRIVATE=github.com/antelman107/net-wait-go go get github.com/antelman107/net-wait-go@v1.3
RUN go install github.com/antelman107/net-wait-go@latest

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN go build -ldflags "-s -w" -o ${BIN_FILE} cmd/social/*

# На выходе тонкий образ
FROM alpine:3.21.3

LABEL SERVICE="social-network"
LABEL MAINTAINERS="usmartpro@gmail.com"

ENV BIN_FILE "/opt/social/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

COPY --from=build /go/bin/net-wait-go /bin/net-wait-go

CMD ${BIN_FILE}
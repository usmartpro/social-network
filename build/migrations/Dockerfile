FROM golang:1.24.1 AS build

ENV CGO_ENABLED=0
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 && \
    go install github.com/antelman107/net-wait-go@latest

FROM alpine:3.21.3
RUN apk add --update --no-cache netcat-openbsd

COPY --from=build /go/bin/goose /bin/goose
COPY --from=build /go/bin/net-wait-go /bin/net-wait-go

WORKDIR /var/migrations

COPY migrations/* /var/migrations/

ENV DRIVER=${GOOSE_DRIVER}
ENV DSN=${GOOSE_DBSTRING}

CMD ["/bin/goose", "up"]
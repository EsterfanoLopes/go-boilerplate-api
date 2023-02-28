FROM golang:1.17.3-alpine3.14 AS builder

RUN apk update && apk add --no-cache ca-certificates gcc musl-dev tzdata && update-ca-certificates && \
        cp /usr/share/zoneinfo/America/Recife /etc/localtime && \
        echo "America/Recife" >  /etc/timezone && \
    apk del tzdata

WORKDIR /go-boilerplate

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags='-w -s -extldflags "-static"' -v -a -o /go/bin/go-boilerplate .

FROM alpine:3.14

COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go-boilerplate/docs/swagger/swagger.json /docs/swagger/swagger.json
COPY --from=builder /go/bin/go-boilerplate /go/bin/go-boilerplate

EXPOSE 9000

ENTRYPOINT ["/go/bin/go-boilerplate"]

FROM golang:1.15 as builder
WORKDIR /go/src/app
COPY . .

RUN make deps && make lint && make build

FROM debian:stretch
COPY --from=builder /go/src/app/bin/eta-service /usr/local/bin/eta-service
COPY --from=builder /go/src/app/values/values_local.yaml ./values/values_local.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["eta-service"]
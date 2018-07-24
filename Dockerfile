FROM golang:1.10.3-alpine3.8 as builder

WORKDIR /go/src/github.com/jsturtevant/azure-k8-metrics-adapter
COPY . .

RUN CGO_ENABLED=0 go build -a -tags netgo -o /adapter github.com/jsturtevant/azure-k8-metrics-adapter

FROM alpine:3.8
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates
    
ENTRYPOINT ["/adapter"]
COPY --from=builder /adapter /
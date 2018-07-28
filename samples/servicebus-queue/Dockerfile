FROM golang:1.10.3-alpine3.8 as builder

WORKDIR /go/src/consumer
RUN apk add -U git
RUN go get -u github.com/Azure/azure-service-bus-go

COPY consumer/ .
RUN CGO_ENABLED=0 go build -a -tags netgo -o /consumer 

FROM alpine:3.8
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates
    
ENTRYPOINT ["/consumer"]
COPY --from=builder /consumer /
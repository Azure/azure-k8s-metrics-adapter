FROM golang:1.10.3-alpine3.8 as builder

WORKDIR /go/src/github.com/Azure/azure-k8s-metrics-adapter
COPY . .

RUN CGO_ENABLED=0 go test $(go list ./... | grep -v -e '/client/' -e '/samples/' -e '/apis/')
RUN CGO_ENABLED=0 go build -a -tags netgo -o /adapter github.com/Azure/azure-k8s-metrics-adapter

FROM alpine:3.8
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates
    
ENTRYPOINT ["/adapter", "--logtostderr=true"]
COPY --from=builder /adapter /

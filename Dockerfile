FROM golang:1.12.12-alpine3.10  as builder

WORKDIR /go/src/github.com/Azure/azure-k8s-metrics-adapter
COPY . .

RUN CGO_ENABLED=0 go test $(go list ./... | grep -v -e '/client/' -e '/samples/' -e '/apis/')
RUN CGO_ENABLED=0 go build -a -tags netgo -o /adapter github.com/Azure/azure-k8s-metrics-adapter

FROM alpine:3.10
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates
    
ENTRYPOINT ["/adapter", "--logtostderr=true"]
COPY --from=builder /adapter /

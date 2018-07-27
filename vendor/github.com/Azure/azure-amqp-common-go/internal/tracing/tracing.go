package tracing

import (
	"context"
	"os"

	"github.com/Azure/azure-amqp-common-go/internal"
	"github.com/opentracing/opentracing-go"
	tag "github.com/opentracing/opentracing-go/ext"
)

// StartSpanFromContext starts a span given a context and applies common library information
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	ApplyComponentInfo(span)
	return span, ctx
}

// ApplyComponentInfo applies amqp common library and network info to the span
func ApplyComponentInfo(span opentracing.Span) {
	tag.Component.Set(span, "github.com/Azure/azure-amqp-common-go")
	span.SetTag("version", common.Version)
	applyNetworkInfo(span)
}

func applyNetworkInfo(span opentracing.Span) {
	hostname, err := os.Hostname()
	if err == nil {
		tag.PeerHostname.Set(span, hostname)
	}
}

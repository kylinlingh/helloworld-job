package sop

import (
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

var (
	grpcTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
)

type ClientNameKey struct{}

func newClientSpanFromRequest(req *http.Request) (opentracing.Span, error) {
	ctx := req.Context()
	tracer := opentracing.GlobalTracer()
	var parentSpanCtx opentracing.SpanContext
	// if parent == nil, this span will be root span.
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentSpanCtx = parent.Context()
	}

	httpClientTag := opentracing.Tag{Key: string(ext.SpanKind), Value: "http client"}
	clientName := ctx.Value(ClientNameKey{}).(string)
	if clientName != "" {
		httpClientTag = opentracing.Tag{Key: string(ext.SpanKind), Value: clientName + " client"}
	}

	opts := []opentracing.StartSpanOption{
		opentracing.ChildOf(parentSpanCtx),
		httpClientTag,
		grpcTag,
	}

	var opName string
	if pc, _, _, ok := runtime.Caller(8); ok {
		funcName := runtime.FuncForPC(pc).Name()
		funcName = funcName[strings.LastIndexAny(funcName, "/.")+1:]
		opName = funcName
	}
	if clientName != "" {
		opName = clientName + "." + opName
	}
	clientSpan := tracer.StartSpan(opName, opts...)
	ext.HTTPMethod.Set(clientSpan, req.Method)
	enEscapeURL, _ := url.QueryUnescape(req.URL.String())
	ext.HTTPUrl.Set(clientSpan, enEscapeURL)
	return clientSpan, tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}

func finishClientSpan(clientSpan opentracing.Span, err error) {
	if err != nil && err != io.EOF {
		ext.Error.Set(clientSpan, true)
		clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}
	clientSpan.Finish()
}

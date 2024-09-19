package sop

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const KB = 1024

func TransportWrapper(rt http.RoundTripper, tracingDisabled bool) http.RoundTripper {
	return &transportWrapper{RoundTripper: rt, TracingDisabled: tracingDisabled}
}

type transportWrapper struct {
	http.RoundTripper
	TracingDisabled bool
}

func (w *transportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var clientSpan opentracing.Span
	if !w.TracingDisabled {
		var err error
		clientSpan, err = newClientSpanFromRequest(req)
		if err != nil {
			ctxzap.Extract(req.Context()).Sugar().With(zap.Error(err)).Warnf("newClientSpanFromRequest failed")
		}
	}
	startTime := time.Now()
	method := req.Method
	url := req.URL.String()

	resp, err := w.RoundTripper.RoundTrip(req)

	if clientSpan != nil {
		finishClientSpan(clientSpan, err)
	}

	latency := time.Since(startTime)
	var statusCode int
	if resp != nil {
		statusCode = resp.StatusCode
	}

	curl := fmt.Sprintf("curl -X %s '%s'", req.Method, req.URL.String())
	for k, values := range req.Header {
		for _, v := range values {
			curl += fmt.Sprintf(" -H '%s'", strings.ReplaceAll(fmt.Sprintf("%s: %s", k, v), `'`, `'\''`))
		}
	}
	if req.Body != nil && req.Body != http.NoBody && req.GetBody != nil {
		if r, _ := req.GetBody(); r != nil {
			if body, _ := ioutil.ReadAll(r); len(body) > 0 {
				s := string(body)
				if len(s) > 10*KB {
					s = s[:10*KB] + "..."
				}
				curl += fmt.Sprintf(" -d '%s'", strings.ReplaceAll(s, `'`, `'\''`))
			}
		}
	}
	var responseHeaders http.Header
	var responseBody []byte
	if resp != nil {
		responseHeaders = resp.Header
		if resp.Body != nil {
			responseBody, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return resp, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(responseBody))
		}
	}

	s := string(responseBody)
	if len(s) > 10*KB {
		s = s[:10*KB] + "..."
	}
	sugar := ctxzap.Extract(req.Context()).Sugar()
	sugar.With(
		zap.Error(err),
		zap.String("curl", base64.StdEncoding.EncodeToString([]byte(curl))),
		zap.Int("response_status", statusCode),
	).Infof("[Call-HTTP] %3d | %13v | %-7s %s", statusCode, latency, method, url)
	sugar.With(
		zap.Error(err),
		zap.String("curl", curl),
		zap.Int("response_status", statusCode),
		zap.Any("response_headers", responseHeaders),
		zap.String("response_body", s),
	).Debugf("[Call-HTTP-Debug] %3d | %13v | %-7s %s", statusCode, latency, method, url)
	return resp, err
}

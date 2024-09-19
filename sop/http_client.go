package sop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpHandledCounter   *prometheus.CounterVec
	httpHandledHistogram *prometheus.HistogramVec
)

func init() {
	httpHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_handled_total",
			Help: "Total number of HTTP request completed by the client.",
		},
		[]string{"http_service", "func", "status_code", "error"},
	)
	httpHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_client_handling_seconds",
			Help:    "Histogram of response latency (seconds) of HTTP until it is finished by the client.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"http_service", "func"},
	)
	_ = prometheus.Register(httpHandledCounter)
	_ = prometheus.Register(httpHandledHistogram)
}

type Config struct {
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host"`
	MaxConnsPerHost     int           `yaml:"max_conns_per_host"`
	IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout"`
	TracingDisabled     bool          `yaml:"tracing_disabled"`
}

func New(name string, conf *Config, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	cli := &Client{
		Client: &http.Client{
			Transport: TransportWrapper(&http.Transport{
				// 读取环境变量里的 HTTP_PROXY 或者 HTTPS_PROXY
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          conf.MaxIdleConns,
				MaxIdleConnsPerHost:   conf.MaxIdleConnsPerHost,
				MaxConnsPerHost:       conf.MaxConnsPerHost,
				IdleConnTimeout:       conf.IdleConnTimeout,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}, conf.TracingDisabled),
			Timeout: timeout,
		},
		name: name,
	}
	return cli
}

type Client struct {
	*http.Client

	name string
}

func (cli *Client) DoRequestJSON(ctx context.Context, method, url string, reqHeader http.Header,
	reqData interface{}, respData ResponseData) (statusCode int, respHeaders http.Header, err error) {
	return doHTTPRequestJSON(ctx, cli.Client, cli.name, method, url, reqHeader, reqData, respData)
}

func (cli *Client) Do(req *http.Request) (resp *http.Response, err error) {
	return do(cli.Client, cli.name, req)
}

func do(cli *http.Client, name string, req *http.Request) (resp *http.Response, err error) {
	startTime := time.Now()
	defer func() {
		var e error
		if err != nil {
			e = err
		} else if resp != nil && resp.StatusCode >= http.StatusBadRequest {
			e = errors.New(resp.Status)
		}
		var statusCode int
		if resp != nil {
			statusCode = resp.StatusCode
		}
		reportPrometheus(name, e, statusCode, startTime)
	}()
	return cli.Do(req.WithContext(context.WithValue(req.Context(), ClientNameKey{}, name)))
}

type ResponseData interface {
	GetError() error
}

func doHTTPRequestJSON(ctx context.Context, cli *http.Client, name, method, url string, reqHeader http.Header,
	reqData interface{}, respData ResponseData) (statusCode int, respHeaders http.Header, err error) {
	startTime := time.Now()
	defer func() {
		reportPrometheus(name, err, statusCode, startTime)
	}()
	ctx = context.WithValue(ctx, ClientNameKey{}, name)
	var reqBody io.Reader
	if reqData != nil {
		b, err := json.Marshal(reqData)
		if err != nil {
			return 0, nil, err
		}
		reqBody = bytes.NewReader(b)
		if reqHeader == nil {
			reqHeader = make(http.Header)
		}
		if len(reqHeader.Values("Content-Type")) == 0 {
			reqHeader.Set("Content-Type", "application/json")
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return 0, nil, err
	}
	if reqHeader != nil {
		req.Header = reqHeader
	}

	resp, err := cli.Do(req)
	if err != nil {
		return 0, nil, err
	}

	if resp != nil {
		respHeaders = resp.Header
	}
	var respBodyBytes []byte
	if resp != nil && resp.Body != nil {
		defer func() { _ = resp.Body.Close() }()
		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, respHeaders, err
		}
	}
	if resp.StatusCode != http.StatusOK {
		_ = json.Unmarshal(respBodyBytes, respData)
		return resp.StatusCode, respHeaders, fmt.Errorf("http status: %s, %s", resp.Status, respBodyBytes)
	}
	err = json.Unmarshal(respBodyBytes, respData)
	if err != nil {
		return resp.StatusCode, respHeaders, fmt.Errorf("unmarshal error: %v", err)
	}
	return resp.StatusCode, respHeaders, respData.GetError()
}

func reportPrometheus(name string, err error, statusCode int, startTime time.Time) {
	latency := time.Since(startTime)
	var errLabel string
	if err != nil {
		errLabel = "1"
	} else {
		errLabel = "0"
	}
	var funcName string
	if pc, _, _, ok := runtime.Caller(4); ok {
		funcName = runtime.FuncForPC(pc).Name()
		funcName = funcName[strings.LastIndexAny(funcName, "/.")+1:]
	}
	httpHandledCounter.WithLabelValues(name, funcName, strconv.Itoa(statusCode), errLabel).Inc()
	httpHandledHistogram.WithLabelValues(name, funcName).Observe(latency.Seconds())
}

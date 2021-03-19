package http

import (
	"bytes"
	"errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type HeaderValue struct {
	Header string
	Value  string
}

func NewHeader(header, value string) HeaderValue {
	return HeaderValue{Header: header, Value: value}
}

type caller interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}

type client struct {
	cli  caller
	head []HeaderValue
}

func NewClient(timeout time.Duration, defaultHeaders ...HeaderValue) *client {
	c := new(client)
	c.head = defaultHeaders

	c.cli = &fasthttp.Client{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	return c
}

func (c *client) Get(url string, headers ...HeaderValue) ([]byte, error) {
	return c.do(c.cli, url, "GET", nil, headers)
}

func (c *client) Post(url string, body []byte, headers ...HeaderValue) ([]byte, error) {
	return c.do(c.cli, url, "POST", body, headers)
}

func (c *client) do(client caller, url, method string, body []byte, headers []HeaderValue) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)

	if body != nil {
		req.SetBody(body)
	}

	for _, h := range c.head {
		req.Header.Set(h.Header, h.Value)
	}
	req.Header.SetMethod(method)
	req.Header.Set("Accept", "application/jsonpb")
	for _, h := range headers {
		req.Header.Set(h.Header, h.Value)
	}

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := client.Do(req, res); err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, errors.New(http.StatusText(res.StatusCode()))
	}

	buf := new(bytes.Buffer)
	buf.Write(res.Body())
	return buf.Bytes(), nil
}

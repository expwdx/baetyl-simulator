package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"baetyl-simulator/errors"
)

var jsonHeaders = map[string]string{"Content-Type": "application/json"}

// Client client of http server
type Client struct {
	ops  *ClientOptions
	http *http.Client
}

// NewClient creates a new http client
func NewClient(ops *ClientOptions) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   ops.Timeout,
			KeepAlive: ops.KeepAlive,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       ops.TLSConfig,
		MaxIdleConns:          ops.MaxIdleConns,
		IdleConnTimeout:       ops.IdleConnTimeout,
		TLSHandshakeTimeout:   ops.TLSHandshakeTimeout,
		ExpectContinueTimeout: ops.ExpectContinueTimeout,
	}
	return &Client{
		ops: ops,
		http: &http.Client{
			Timeout:   ops.Timeout,
			Transport: transport,
		},
	}
}

// Call calls the function via HTTP POST
func (c *Client) Call(function string, payload []byte) ([]byte, error) {
	r, err := c.PostURL(function, bytes.NewBuffer(payload), jsonHeaders)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return HandleResponse(r)
}

// PostJSON post data with json content type
func (c *Client) PostJSON(url string, payload []byte, headers ...map[string]string) ([]byte, error) {
	headers = append(headers, jsonHeaders)
	r, err := c.PostURL(url, bytes.NewBuffer(payload), headers...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return HandleResponse(r)
}

// GetJSON get data with json content type
func (c *Client) GetJSON(url string, headers ...map[string]string) ([]byte, error) {
	headers = append(headers, jsonHeaders)
	r, err := c.GetURL(url, headers...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return HandleResponse(r)
}

func (c *Client) GetURL(url string, header ...map[string]string) (*http.Response, error) {
	return c.SendUrl("GET", url, nil, header...)
}

func (c *Client) PostURL(url string, body io.Reader, header ...map[string]string) (*http.Response, error) {
	return c.SendUrl("POST", url, body, header...)
}

func (c *Client) SendUrl(method, url string, body io.Reader, header ...map[string]string) (*http.Response, error) {
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("%s/%s", c.ops.Address, url)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for _, v := range header {
		for kk, vv := range v {
			req.Header.Set(kk, vv)
		}
	}
	r, err := c.http.Do(req)
	return r, errors.Trace(err)
}

// HandleResponse handles response
func HandleResponse(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		msg := strings.TrimRight(string(data), "\n")
		if msg == "" {
			msg = r.Status
		}
		err = errors.Errorf("[%d] %s", r.StatusCode, msg)
	}
	return data, errors.Trace(err)
}

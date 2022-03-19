package baetyl

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"baetyl-simulator/config"
	"baetyl-simulator/constants"
	"baetyl-simulator/middleware/log"
)

const CONTENT_TYPE = "application/json; charset=utf-8"

type PathField map[string]string
type QueryParam map[string]string

type BaetylHttpClient interface {
	Get(path string, keyParams PathField, params QueryParam) (*http.Response, error)
	Post(path string, keyParams PathField, body interface{}) (*http.Response, error)
	Put(path string, keyParams PathField, body interface{}) (*http.Response, error)
	Patch(path string, keyParams PathField, body interface{}) (*http.Response, error)
	Delete(path string, keyParams PathField) (*http.Response, error)

	Read(resp *http.Response, view interface{}) error
	ReadMap(resp *http.Response) (map[string]interface{}, error)
}

// baetylHttpClient baetylHttpClient api rpc
type baetylHttpClient struct {
	client *http.Client
	cfg    *config.ServerConfig
}

func NewBaetylClient(ctx context.Context, config *config.ServerConfig) (BaetylHttpClient, error) {
	c := &baetylHttpClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		cfg: config,
	}
	if config.Schema == "https" {
		if ctx != nil {
			name := ctx.Value(constants.NODE_NAME).(string)
			certKeyFile := fmt.Sprintf(constants.NODE_CERT_KEY, name)
			certPemFile := fmt.Sprintf(constants.NODE_CERT_CERT, name)

			// x509.Certificate
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(certPemFile)
			if err != nil {
				return nil, err
			}
			pool.AppendCertsFromPEM(caCrt)

			cliCrt, err := tls.LoadX509KeyPair(certPemFile, certKeyFile)
			if err != nil {
				return nil, err
			}

			//    tr := &http2.Transport{  // http2协议
			c.client.Transport = &http.Transport{ // http1.1协议
				TLSClientConfig: &tls.Config{
					RootCAs:            pool,
					Certificates:       []tls.Certificate{cliCrt},
					InsecureSkipVerify: true,
				},
			}
		} else {
			c.client.Transport = &http.Transport{ // http1.1协议
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}

	}

	return c, nil
}

func (c *baetylHttpClient) Get(path string, keyParams PathField, params QueryParam) (*http.Response, error) {
	u := NewUrl(c.cfg, Path(path), PathFields(keyParams), QueryParams(params))
	return c.client.Get(u.String())
}

func (c *baetylHttpClient) Post(path string, keyParams PathField, body interface{}) (*http.Response, error) {
	u := NewUrl(c.cfg, Path(path), PathFields(keyParams))
	var rb []byte

	switch body.(type) {
	case []byte:
		rb = body.([]byte)
	case string:
		rb = body.([]byte)
	default:
		var err error
		rb, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	return c.client.Post(u.String(), CONTENT_TYPE, bytes.NewBuffer(rb))
}

func (c *baetylHttpClient) Put(path string, keyParams PathField, body interface{}) (*http.Response, error) {
	u := NewUrl(c.cfg, Path(path), PathFields(keyParams))
	rb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", CONTENT_TYPE)
	return c.client.Do(req)
}

func (c *baetylHttpClient) Patch(path string, keyParams PathField, body interface{}) (*http.Response, error) {
	u := NewUrl(c.cfg, Path(path), PathFields(keyParams))
	rb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, u.String(), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", CONTENT_TYPE)
	return  c.client.Do(req)
}

func (c *baetylHttpClient) Delete(path string, keyParams PathField) (*http.Response, error) {
	u := NewUrl(c.cfg, Path(path), PathFields(keyParams))
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", CONTENT_TYPE)
	return  c.client.Do(req)
}

// Close 关闭 io reader
func (c *baetylHttpClient) Close(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.L().Error("client reader close fail", log.Any("error", log.Code(err)))
	}
}

func (c *baetylHttpClient) Read(resp *http.Response, view interface{}) error {
	//defer c.Close(resp.Body)

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBytes, view)
	if err != nil {
		return err
	}

	return err
}

func (c *baetylHttpClient) ReadMap(resp *http.Response) (map[string]interface{}, error) {
	//defer c.Close(resp.Body)

	res := make(map[string]interface{})
	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, err
}


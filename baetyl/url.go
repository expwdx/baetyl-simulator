package baetyl

import (
	"baetyl-simulator/config"
	"net/url"
	"strings"
)

type Option func(*Url)

type Url struct {
	*url.URL
}

func NewUrl(config *config.ServerConfig, options ...func(*Url)) *Url {
	hurl := &Url{
		&url.URL{
			Scheme: config.Schema,
			Host: config.Host,
			Path: config.ApiVer,
		},
	}

	for _, option := range options {
		option(hurl)
	}

	return hurl
}

// Path 请求路径
func Path(path string) Option {
	return func(u *Url) {
		u.Path += path
	}
}

// PathFields 格式化url的path，将路径参数渲染到path
func PathFields(Fields PathField) Option {
	return func(u *Url) {
		if Fields != nil {
			kvs := make([]string, 2 * len(Fields))
			for key, val := range Fields {
				kvs = append(kvs, ":"+key, val)
			}

			r := strings.NewReplacer(kvs...)
			u.Path = r.Replace(u.Path)
		}
	}
}

// QueryParams 组装请求参数
func QueryParams(params QueryParam) Option {
	return func(u *Url) {
		if params != nil {
			q := make(url.Values, len(params))
			for key, val := range params {
				q.Set(key, val)
			}

			u.RawQuery = q.Encode()
		}
	}
}

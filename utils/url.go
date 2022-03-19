package utils

import (
	"net/url"
	"strings"

	"baetyl-simulator/errors"
)

// ParseURL parses a url string
func ParseURL(addr string) (*url.URL, error) {
	if strings.HasPrefix(addr, "unix://") {
		parts := strings.SplitN(addr, "://", 2)
		return &url.URL{
			Scheme: parts[0],
			Host:   parts[1],
		}, nil
	}
	res, err := url.Parse(addr)
	return res, errors.Trace(err)
}

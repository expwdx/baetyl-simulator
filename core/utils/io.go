package utils

import (
	"io"

	"baetyl-simulator/middleware/log"
)

func CloseReader(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.L().Error("client reader close fail", log.Any("error", err))
	}
}
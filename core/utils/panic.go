package utils

import (
	"runtime/debug"

	"baetyl-simulator/errors"
	"baetyl-simulator/middleware/log"
)

func HandlePanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = errors.New("panic error")
		}
		log.L().Error("handle a panic",  log.Error(err), log.Any("panic", string(debug.Stack())))
	}
}

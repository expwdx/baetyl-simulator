package utils

import (
    "baetyl-simulator/errors"
    log2 "baetyl-simulator/middleware/log"
    "io/ioutil"
)

func ReadFile(fileName string) string {
    b, err := ioutil.ReadFile(fileName) // just pass the file name
    if err != nil {
        log2.L().Error("read file fail", log2.Any("file", fileName), log2.Any("error", errors.Trace(err)))
    }

    str := string(b)

    return str
}

package common

import (
    "baetyl-simulator/config"
    log2 "baetyl-simulator/middleware/log"
)

var Log *log2.Logger

func InitLog(config *config.Config) error {
    var lfs []log2.Field

    lfs = append(lfs, log2.Any("app", "baetyl-simulate"))
    Log = log2.With(lfs...)
    Log.Info("to load config file", log2.Any("file", "ctc/conf.yaml"))
    _log, err := log2.Init(config.Logger, lfs...)
    if err != nil {
        return err
    }

    Log = _log

    return nil
}
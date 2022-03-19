package core

import (
    "baetyl-simulator/core/mock"
    "baetyl-simulator/middleware/log"
    "github.com/pkg/errors"
    "github.com/spf13/viper"

    "baetyl-simulator/config"
)

func Clean() error {
    log.L().Info("start initializing simulator...")

    var cfg config.Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.L().Error("load config fail before initializing data.")
        return err
    }

    log.L().Info("load config success.", log.Any("config", cfg))

    mockService, err := mock.NewMockService(&cfg)
    if err != nil {
        return errors.Wrap(err, "error")
    }

    err = mockService.CleanData()
    if err != nil {
        return errors.Wrap(err, "simulator init fail")
    }

    return nil
}

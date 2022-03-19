package core

import (
    "baetyl-simulator/core/mock"
    log2 "baetyl-simulator/middleware/log"
    "github.com/pkg/errors"
    "github.com/spf13/viper"

    "baetyl-simulator/config"
)

func Initialize() error {
    log2.L().Info("start initializing simulator...")

    var cfg config.Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log2.L().Error("load config fail before initializing data.")
        return err
    }

    log2.L().Info("load config success.", log2.Any("config", cfg))

    mockService, err := mock.NewMockService(&cfg)
    if err != nil {
        return errors.Wrap(err, "error")
    }

    err = mockService.InitData()
    if err != nil {
        return errors.Wrap(err, "simulator init fail")
    }

    return nil
}

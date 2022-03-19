package core

import (
    "baetyl-simulator/common"
    "baetyl-simulator/core/cloud"
    "baetyl-simulator/core/engine"
    "context"
    "fmt"
    "github.com/gin-contrib/cache/persistence"
    "github.com/pkg/errors"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/spf13/viper"

    "baetyl-simulator/config"
    "baetyl-simulator/middleware/log"
)

func waitChan() <-chan os.Signal {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
    signal.Ignore(syscall.SIGPIPE)

    return sig
}

func Wait() {
    <-waitChan()
}

func Run() error {
    log.L().Info("start initializing simulator...")

    var cfg config.Config
    if err := viper.Unmarshal(&cfg); err != nil {
        log.L().Error("load config fail before initializing data.")
        return err
    }

    common.Cache = persistence.NewInMemoryStore(cfg.Test.PlanTime)

    log.L().Info("load config success.", log.Any("config", cfg))
    ctx, stop := context.WithCancel(context.Background())
    defer stop()

    a, err := cloud.NewAdmin(&cfg)
    if err != nil {
        return err
    }

    go a.Start(ctx)

    for i := cfg.Test.NodeStartNo; i < cfg.Test.NodeCount; i++ {
        nodeName := fmt.Sprintf("%s-%d", cfg.Mock.NodeNamePrefix, i)
        e, err := engine.NewEngine(ctx, nodeName, &cfg)
        if err != nil {
            return errors.Wrap(err, "new engine fail")
        }

        go e.Start(ctx)

        time.Sleep(10 * time.Millisecond)
        log.L().Info("startup node", log.Any("node", i))
    }

    t := cfg.Test.PlanTime
    switch t > 0 {
    case true:
        time.Sleep(t)
    default:
        Wait()
    }

    return nil
}

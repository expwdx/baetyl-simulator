package common

import (
    log2 "baetyl-simulator/middleware/log"
    "os"
    "os/signal"
    "runtime/debug"
    "syscall"
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

func Run(handle func() error) {
    defer func() {
        if r := recover(); r != nil {
            Log.Error("service is stopped with panic", log2.Any("panic", r), log2.Any("stack", string(debug.Stack())))
        }
    }()

    err := handle()
    if err != nil {
        log2.L().Error("service has stopped with error", log2.Error(err))
    } else {
        log2.L().Info("service has stopped")
    }
}
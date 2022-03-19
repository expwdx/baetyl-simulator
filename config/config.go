package config

import (
    log2 "baetyl-simulator/middleware/log"
    "time"
)

// Config config
type Config struct {
    Test        TestConfig      `yaml:"test" json:"test"`
    User        UserConfig      `yaml:"user" json:"user"`
    Engine      EngineConfig    `yaml:"engine" json:"engine"`
    Logger      log2.Config     `yaml:"logger" json:"logger"`
    Lock        Lock            `yaml:"lock" json:"lock"`
    Cloud       CloudConfig     `yaml:"cloud" json:"cloud" default:"{}"`
    Mock        MockConfig      `yaml:"mock" json:"mock"`
    Template    Template        `yaml:"template"`
}

type TestConfig struct {
    PlanTime            time.Duration `yaml:"planTime" json:"planTime"`
    NodeCount           int           `yaml:"nodeCount" json:"nodeCount"`
    NodeStartNo         int           `yaml:"nodeStartNo" default:"0"`
}

type UserConfig struct {
    Deploy struct {
        Interval time.Duration `yaml:"interval" json:"interval" default:"10s"`
    } `yaml:"deploy" json:"deploy"`
    Read struct {
        Interval time.Duration `yaml:"interval" json:"interval" default:"20s"`
    } `yaml:"read" json:"read"`
}

type EngineConfig struct {
    Report struct {
        Interval time.Duration `yaml:"interval" json:"interval" default:"10s"`
    } `yaml:"report" json:"report"`
    Desire struct {
        Interval time.Duration `yaml:"interval" json:"interval" default:"20s"`
    } `yaml:"desire" json:"desire"`
}

type Lock struct {
    ExpireTime int64 `yaml:"expireTime" json:"expireTime" default:"5" unit:"second"`
}

type CloudConfig struct {
    Admin ServerConfig `yaml:"admin" json:"admin" default:"{}"`
    Init  ServerConfig `yaml:"init" json:"init" default:"{}"`
    Sync  ServerConfig `yaml:"sync" json:"sync" default:"{}"`
}

type ServerConfig struct {
    Schema  string        `yaml:"schema" default:"http"`
    Host    string        `yaml:"host" validate:"nonzero"`
    ApiVer  string        `yaml:"apiVer" default:"v1"`
    Timeout time.Duration `yaml:"timeout,omitempty" default:"30s"`
}

type MockConfig struct {
    Namespace           string                  `yaml:"namespace" default:"baetyl-cloud"`
    NodeNamePrefix      string                  `yaml:"nodeNamePrefix" default:"simulator_node_"`
    NodeCount           int                     `yaml:"nodeCount" default:"100"`
    NodeStartNo         int                     `yaml:"nodeStartNo" default:"0"`
    NodeLabels          map[string]string       `yaml:"nodeLabels"`
    AppName             string                  `yaml:"appName"`
}

type Template struct {
    Path    string      `yaml:"path"`
}

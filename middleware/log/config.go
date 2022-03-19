package log

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"baetyl-simulator/errors"
)

// Config for logging
type Config struct {
	Level       string                   `yaml:"level" json:"level" default:"info" validate:"regexp=^(fatal|panic|error|warn|info|debug)$"`
	Encoding    string                   `yaml:"encoding" json:"encoding" default:"json" validate:"regexp=^(json|console)$"`
	Filename    string                   `yaml:"filename" json:"filename"`
	Compress    bool                     `yaml:"compress" json:"compress"`
	MaxAge      int                      `yaml:"maxAge" json:"maxAge" default:"15" validate:"min=1"`   // days
	MaxSize     int                      `yaml:"maxSize" json:"maxSize" default:"50" validate:"min=1"` // MB
	MaxBackups  int                      `yaml:"maxBackups" json:"maxBackups" default:"15" validate:"min=1"`
	EncodeTime  string                   `yaml:"encodeTime" json:"encodeTime"`   // time format, like [2006/01/02 15:04:05 UTC]
	EncodeLevel string                   `yaml:"encodeLevel" json:"encodeLevel"` // symbols surround level, like [level]
	EnableKafka bool                     `yaml:"enableKafka" json:"enableKafka" default:"false"` //enable kafka log backend
	KafkaLoggerConfig *KafkaLoggerConfig `yaml:"kafkaLoggerConfig,omitempty" json:"kafkaLoggerConfig,omitempty"`
}

type KafkaLoggerConfig struct {
	Hosts       []string   `yaml:"hosts" json:"hosts" default:"[]"`
	Topic       string     `yaml:"topic" json:"topic" default:""`
}

type KafkaLogInfo struct {
	TraceId       string        `json:"trace_id" default:""`
	CreateTime    string        `json:"create_time"`
	Method        string        `json:"method" default:""`
	SpanId        string        `json:"span_id" default:""`
	Level         string        `json:"level"`
	ErrorDetail   string        `json:"error_detail,omitempty"`
	MsgInfo       string        `json:"msginfo"`
	ThreadId      string        `json:"threadId" default:""`
	Topic         string        `json:"topic"`
	Class         string        `json:"class" default:""`
	K8sPodName    string        `json:"k8s_pod_name" default:""`
}

func (c *Config) String() string {
	res := fmt.Sprintf(
		"level=%s&encoding=%s&filename=%s&compress=%t&maxAge=%d&maxSize=%d&maxBackups=%d&enableKafka=%t",
		c.Level,
		c.Encoding,
		base64.URLEncoding.EncodeToString([]byte(c.Filename)),
		c.Compress,
		c.MaxAge,
		c.MaxSize,
		c.MaxBackups,
		c.EnableKafka)

	if c.EnableKafka == true {
		res = res + fmt.Sprintf("&hosts=%s&topic=%s",
			strings.Join(c.KafkaLoggerConfig.Hosts, ","),
			c.KafkaLoggerConfig.Topic)
	}

	return res
}

// FromURL creates config from url
func FromURL(u *url.URL) (*Config, error) {
	args := u.Query()
	c := new(Config)
	c.Level = args.Get("level")
	c.Encoding = args.Get("encoding")
	filename, err := base64.URLEncoding.DecodeString(args.Get("filename"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.Filename = string(filename)
	c.Compress, err = strconv.ParseBool(args.Get("compress"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.MaxAge, err = strconv.Atoi(args.Get("maxAge"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.MaxSize, err = strconv.Atoi(args.Get("maxSize"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.MaxBackups, err = strconv.Atoi(args.Get("maxBackups"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	c.EnableKafka, err = strconv.ParseBool(args.Get("enableKafka"))
	if err != nil {
		return nil, errors.Trace(err)
	}

	if c.EnableKafka == true {
		kcfg := &KafkaLoggerConfig{
			Hosts: strings.Split(args.Get("hosts"), ","),
			Topic: args.Get("topic"),
		}
		c.KafkaLoggerConfig = kcfg
	}

	return c, nil
}

type LogEntry struct {
	Level       string       `json:"level,omitempty"`
	Time        string       `json:"ts,omitempty"`
	LoggerName  string       `json:"logger,omitempty"`
	Message     string       `json:"msg,omitempty"`
	Caller      string       `json:"caller,omitempty"`
	Stack       string       `json:"stacktrace,omitempty"`
	Function    string       `json:"function,omitempty"`
	Method      string       `json:"method,omitempty"`
	Url         string       `json:"url,omitempty"`
	Pwd         string       `json:"pwd,omitempty"`
	Args        []string     `json:"args,omitempty"`
}

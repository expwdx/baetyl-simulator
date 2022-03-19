package log

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	kafkaSinkInsts = map[string]KafkaSink{}
)

const kafkaDefaultTopic = "normal-baetyl-cloud"

type KafkaSink struct {
	kafkaProducer sarama.SyncProducer
	isAsync       bool
	topic         string
}

func getKafkaSink(brokers []string, topic string, config *sarama.Config) KafkaSink {
	producerInst, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	kafkaSinkInst := KafkaSink{
		kafkaProducer: producerInst,
		topic:         topic,
	}
	return kafkaSinkInst
}

// InitKafkaSink  create kafka sink instance
func InitKafkaSink(u *url.URL) (zap.Sink, error) {
	topic := kafkaDefaultTopic
	var hosts []string
	if t := u.Query().Get("topic"); len(t) > 0 {
		topic = t
	}
	if h := u.Query().Get("hosts"); len(h) > 0 {
		hosts = strings.Split(h, ",")
	}
	brokers := hosts
	instKey := strings.Join(brokers, ",")
	if v, ok := kafkaSinkInsts[instKey]; ok {
		return v, nil
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	if ack := u.Query().Get("acks"); len(ack) > 0 {
		if iack, err := strconv.Atoi(ack); err == nil {
			config.Producer.RequiredAcks = sarama.RequiredAcks(iack)
		} else {
			log.Printf("kafka producer acks value '%s' invalid  use default value %d\n", ack, config.Producer.RequiredAcks)
		}
	}
	if retries := u.Query().Get("retries"); len(retries) > 0 {
		if iretries, err := strconv.Atoi(retries); err == nil {
			config.Producer.Retry.Max = iretries
		} else {
			log.Printf("kafka producer retries value '%s' invalid  use default value %d\n", retries, config.Producer.Retry.Max)
		}
	}
	kafkaSinkInsts[instKey] = getKafkaSink(brokers, topic, config)
	return kafkaSinkInsts[instKey], nil
}

// Close implement zap.Sink func Close
func (s KafkaSink) Close() error {
	return nil
}

// Write implement zap.Sink func Write
func (s KafkaSink) Write(b []byte) (n int, err error) {
	var multiErr error

	logEntry := new(LogEntry)
	if err := json.Unmarshal(b, &logEntry); err != nil {
		multiErr = multierr.Append(multiErr, err)
	}

	k8sPodName := os.Getenv("HOSTNAME")
	msgInfo := &KafkaLogInfo{
		CreateTime: logEntry.Time[:23],
		Level:      strings.ToUpper(logEntry.Level),
		MsgInfo:    string(b),
		Topic:      s.topic,
		K8sPodName: k8sPodName,
		ThreadId:   "",
		Method: logEntry.Method,
	}

	msg, err := json.Marshal(msgInfo)
	if err != nil {
		log.Printf(err.Error())
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.StringEncoder(time.Now().String()),
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		multiErr = multierr.Append(multiErr, err)
	}
	return len(b), multiErr
}

// Sync implement zap.Sink func Sync
func (s KafkaSink) Sync() error {
	return nil
}


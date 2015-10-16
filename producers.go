package main

import (
	"log"

	"github.com/bitly/go-nsq"
)

var (
	producer *nsq.Producer
)

func newProducer() *nsq.Producer {
	nsqConfig.UserAgent = "cqrsnuProducer"

	producer, err := nsq.NewProducer(nsqdTCPAddr, nsqConfig)
	if err != nil {
		log.Fatalf("error creating nsq.Producer: %s", err)
	}
	producer.SetLogger(nil, 0)

	if err = producer.Ping(); err != nil {
		log.Fatalf("error pinging nsqd: %s\n", err)
	}

	return producer
}

func Send(topic string, message interface{}) {
	if producer == nil {
		producer = newProducer()
	}

	if err := producer.Publish(topic, Serialize(message)); err != nil {
		log.Fatalf("Send %s error: %s\n", topic, err)
	}
}

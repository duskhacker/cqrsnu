package main

import (
	"log"

	"github.com/bitly/go-nsq"
)

var (
	consumers []*nsq.Consumer
)

func InitConsumers() {
	newConsumer(openTab, openTab+"Consumer", OpenTabHandler)
	newConsumer(placeOrder, placeOrder+"Consumer", PlaceOrderHandler)
	newConsumer(markDrinksServed, markDrinksServed+"Consumer", MarkDrinksServedHandler)
	newConsumer(markFoodPrepared, markFoodPrepared+"Consumer", MarkFoodPreparedHandler)
	newConsumer(drinksServed, drinksServed+"Consumer", DrinksServedHandler)
	newConsumer(closeTab, closeTab+"Consumer", CloseTabHandler)
}

func newConsumer(topic, channel string, handler func(*nsq.Message) error) *nsq.Consumer {
	nsqConfig.UserAgent = channel

	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		log.Fatalf("%s:%s; NewConsumer: %s", topic, channel, err)
	}
	consumer.SetLogger(nil, 0)

	consumer.AddHandler(nsq.HandlerFunc(nsq.HandlerFunc(handler)))

	if connectToNSQD {
		if err = consumer.ConnectToNSQD(nsqdTCPAddr); err != nil {
			log.Fatalf("%s:%s; ConnectToNSQLookupds: %s", topic, channel, err)
		}

	} else {
		if err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs); err != nil {
			log.Fatalf("%s:%s; ConnectToNSQLookupds: %s", topic, channel, err)
		}
	}

	consumers = append(consumers, consumer)
	return consumer
}

func StopAllConsumers() {
	for _, consumer := range consumers {
		consumer.Stop()
	}
}

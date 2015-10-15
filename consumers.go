package main

import (
	"log"

	"github.com/bitly/go-nsq"
)

var (
	openTabConsumer          *nsq.Consumer
	tapOpenedConsumer        *nsq.Consumer
	placeOrderConsumer       *nsq.Consumer
	markDrinksServedConsumer *nsq.Consumer
	drinksServedConsumer     *nsq.Consumer
)

func initConsumers() {
	openTabConsumer = newConsumer(openTab, openTab+"Consumer", OpenTabHandler)
	placeOrderConsumer = newConsumer(placeOrder, placeOrder+"Consumer", PlaceOrderHandler)
	markDrinksServedConsumer = newConsumer(markDrinksServed, markDrinksServed+"Consumer", MarkDrinksServedHandler)
	drinksServedConsumer = newConsumer(drinksServed, drinksServed+"Consumer", DrinksServedHandler)
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
	return consumer
}

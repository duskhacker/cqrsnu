package main

import (
	"fmt"
	"log"

	"github.com/bitly/go-nsq"
)

var (
	openTabProducer          *nsq.Producer
	tabOpenedProducer        *nsq.Producer
	placeOrderProducer       *nsq.Producer
	foodOrderedProducer      *nsq.Producer
	drinksOrderedProducer    *nsq.Producer
	drinksServedProducer     *nsq.Producer
	markDrinksServedProducer *nsq.Producer
	exceptionProducer        *nsq.Producer
)

func newProducer(channel, topic string) *nsq.Producer {
	nsqConfig.UserAgent = fmt.Sprintf("%sProducer", topic)

	producer, err := nsq.NewProducer(nsqdTCPAddr, nsqConfig)
	if err != nil {
		log.Fatalf("%s;%s failed to create nsq.Producer - %s", topic, channel, err)
	}
	producer.SetLogger(nil, 0)

	if err = producer.Ping(); err != nil {
		log.Fatalf("%s:%s; unable to ping nsqd: %s\n", topic, channel, err)
	}

	return producer
}

func Send(topic string, message interface{}) {
	switch topic {
	case openTab:
		if openTabProducer == nil {
			openTabProducer = newProducer(openTab, openTab+"Producer")
		}
		if err := openTabProducer.Publish(openTab, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", openTab, err)
		}
	case tabOpened:
		if tabOpenedProducer == nil {
			tabOpenedProducer = newProducer(tabOpened, tabOpened+"Producer")
		}
		if err := tabOpenedProducer.Publish(tabOpened, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", tabOpened, err)
		}
	case placeOrder:
		if placeOrderProducer == nil {
			placeOrderProducer = newProducer(placeOrder, placeOrder+"Producer")
		}
		if err := placeOrderProducer.Publish(placeOrder, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", placeOrder, err)
		}
	case foodOrdered:
		if foodOrderedProducer == nil {
			foodOrderedProducer = newProducer(foodOrdered, foodOrdered+"Producer")
		}
		if err := foodOrderedProducer.Publish(foodOrdered, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", foodOrdered, err)
		}
	case drinksOrdered:
		if drinksOrderedProducer == nil {
			drinksOrderedProducer = newProducer(drinksOrdered, drinksOrdered+"Producer")
		}
		if err := drinksOrderedProducer.Publish(drinksOrdered, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", drinksOrdered, err)
		}
	case drinksServed:
		if drinksServedProducer == nil {
			drinksServedProducer = newProducer(drinksServed, drinksServed+"Producer")
		}
		if err := drinksServedProducer.Publish(drinksServed, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", drinksServed, err)
		}
	case markDrinksServed:
		if markDrinksServedProducer == nil {
			markDrinksServedProducer = newProducer(markDrinksServed, markDrinksServed+"Producer")
		}
		if err := markDrinksServedProducer.Publish(markDrinksServed, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", markDrinksServed, err)
		}
	case exception:
		if exceptionProducer == nil {
			exceptionProducer = newProducer(exception, exception+"Producer")
		}
		if err := exceptionProducer.Publish(exception, Serialize(message)); err != nil {
			log.Fatalf("Send %s error: %s\n", exception, err)
		}
	default:
		log.Fatalf("Unknown topic: %s\n", topic)

	}
}

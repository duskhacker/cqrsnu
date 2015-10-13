package main

import (
	"fmt"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/internal/app"
)

//public class TabAggregate : Aggregate
//{
//}

const (
	maxConcurrentHttpRequests = 5
	tabOpenedTopic            = "TabOpened"
	openTabTopic              = "OpenTab"
)

type Tab struct {
	Guid uuid.UUID
}

var (
	lookupdHTTPAddrs = app.StringArray{}
	nsqdTCPAddr      = "localhost:4150"

	nsqConfig         = nsq.NewConfig()
	tabOpenedProducer *nsq.Producer
	openTabConsumer   *nsq.Consumer
)

func OpenTabHandler(msg *nsq.Message) error {
	ot := OpenTab{}.FromJson(msg.Body)
	tabOpenedProducer.Publish("TabOpened", NewTabOpened(ot.Guid, ot.TableNumber, ot.WaitStaff).ToJson())
	return nil
}

func initHandlers() {
	tabOpenedProducer = newProducer(tabOpenedTopic, tabOpenedTopic+"Producer")
	openTabConsumer = newConsumer(openTabTopic, openTabTopic+"Consumer", OpenTabHandler)
}

func newConsumer(topic, channel string, handler func(*nsq.Message) error) *nsq.Consumer {
	nsqConfig.UserAgent = channel

	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		log.Fatalf("%s:%s; NewConsumer: %s", topic, channel, err)
	}

	consumer.AddHandler(nsq.HandlerFunc(nsq.HandlerFunc(handler)))

	if err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs); err != nil {
		log.Fatalf("%s:%s; ConnectToNSQLookupds: %s", topic, channel, err)
	}
	return consumer
}

func newProducer(channel, topic string) *nsq.Producer {
	nsqConfig.UserAgent = fmt.Sprintf("%sProducer", topic)

	producer, err := nsq.NewProducer(nsqdTCPAddr, nsqConfig)
	if err != nil {
		log.Fatalf("%s;%s failed to create nsq.Producer - %s", topic, channel, err)
	}

	if err = producer.Ping(); err != nil {
		log.Fatalf("%s:%s; unable to ping nsqd: %s\n", topic, channel, err)
	}

	return producer
}

func main() {
	nsqConfig.MaxInFlight = maxConcurrentHttpRequests
}

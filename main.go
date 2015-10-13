package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/bitly/nsq/internal/app"
)

//public class TabAggregate : Aggregate
//{
//}

const (
	Topic                     = "Events"
	channel                   = "Cafe"
	maxConcurrentHttpRequests = 5
)

type Tab struct {
	Guid uuid.UUID
}

var (
	lookupdHTTPAddrs = app.StringArray{}
	Tabs             = []Tab{}
	TabOpenedStore   = []TabOpened{}
)

//func HandleMessage(msg *nsq.Message) error {
//	EventStore = append(O)
//
//}

//func doSomething(msg interface{}) {
//	fmt.Printf("%#v\n", msg)
//}

//func NsqSetup() {
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
//
//	cfg := nsq.NewConfig()
//	cfg.UserAgent = channel
//	cfg.MaxInFlight = maxConcurrentHttpRequests
//
//	consumer, err := nsq.NewConsumer(Topic, channel, cfg)
//	if err != nil {
//		log.Fatalf("[rir] %s", err)
//	}
//
//	consumer.AddHandler(nsq.HandlerFunc(HandleMessage))
//
//	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
//	if err != nil {
//		log.Fatalf("[rir] %s", err)
//	}
//
//	for {
//		select {
//		case <-consumer.StopChan:
//			return
//		case <-sigChan:
//			consumer.Stop()
//		}
//	}
//
//}

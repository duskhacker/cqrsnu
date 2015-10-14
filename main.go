package main

import (
	"encoding/json"
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
	maxConcurrentHttpRequests = 100
	openTab                   = "OpenTab"
	tabOpened                 = "TabOpened"
	placeOrder                = "PlaceOrder"
	foodOrdered               = "FoodOrdered"
	drinksOrdered             = "DrinksOrdered"
	exception                 = "Exception"
)

var (
	lookupdHTTPAddrs = app.StringArray{}
	nsqdTCPAddr      = "localhost:4150"

	nsqConfig             = nsq.NewConfig()
	openTabProducer       *nsq.Producer
	tabOpenedProducer     *nsq.Producer
	placeOrderProducer    *nsq.Producer
	foodOrderedProducer   *nsq.Producer
	drinksOrderedProducer *nsq.Producer
	exceptionProducer     *nsq.Producer

	openTabConsumer    *nsq.Consumer
	tapOpenedConsumer  *nsq.Consumer
	placeOrderConsumer *nsq.Consumer

	Tabs = make(map[string]Tab)
)

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

func OpenTabHandler(msg *nsq.Message) error {
	ot := OpenTab{}.FromJson(msg.Body)
	Tabs[ot.ID.String()] = Tab{TableNumber: ot.TableNumber, WaitStaff: ot.WaitStaff}
	Send(tabOpened, NewTabOpened(ot.ID, ot.TableNumber, ot.WaitStaff))
	return nil
}

func Serialize(object interface{}) []byte {
	o, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("error marshaling %#v: %s", object, err)
	}
	return o
}

func PlaceOrderHandler(msg *nsq.Message) error {
	order := new(PlaceOrder).FromJson(msg.Body)
	tab, ok := Tabs[order.ID.String()]
	if !ok {
		Send(exception, NewCommandException(uuid.NewRandom(), "TabNotOpen", "Cannot Place order without open Tab"))
		return nil
	}

	var (
		foodItems  []OrderedItem
		drinkItems []OrderedItem
	)

	for _, item := range order.Items {
		if item.IsDrink {
			drinkItems = append(drinkItems, item)
		} else {
			foodItems = append(foodItems, item)
		}
	}

	if len(foodItems) > 0 {
		Send(foodOrdered, NewFoodOrdered(order.ID, foodItems))
	}

	if len(drinkItems) > 0 {
		Send(drinksOrdered, NewDrinksOrdered(order.ID, drinkItems))
	}

	tab.Items = append(tab.Items, order.Items...)

	return nil
}

func initConsumers() {
	openTabConsumer = newConsumer(openTab, openTab+"Consumer", OpenTabHandler)
	placeOrderConsumer = newConsumer(placeOrder, placeOrder+"Consumer", PlaceOrderHandler)
}

func newConsumer(topic, channel string, handler func(*nsq.Message) error) *nsq.Consumer {
	nsqConfig.UserAgent = channel

	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		log.Fatalf("%s:%s; NewConsumer: %s", topic, channel, err)
	}
	consumer.SetLogger(nil, 0)

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
	producer.SetLogger(nil, 0)

	if err = producer.Ping(); err != nil {
		log.Fatalf("%s:%s; unable to ping nsqd: %s\n", topic, channel, err)
	}

	return producer
}

func main() {
	nsqConfig.MaxInFlight = maxConcurrentHttpRequests
	initConsumers()
}

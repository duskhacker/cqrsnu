package main

import (
	"encoding/json"
	"log"

	"github.com/bitly/go-nsq"
	"github.com/bitly/nsq/internal/app"
)

const (
	maxConcurrentHttpRequests = 100
)

var (
	lookupdHTTPAddrs = app.StringArray{}
	nsqdTCPAddr      = "localhost:4150"

	nsqConfig = nsq.NewConfig()
	Tabs      = make(map[string]Tab)
)

func Serialize(object interface{}) []byte {
	o, err := json.Marshal(object)
	if err != nil {
		log.Fatalf("error marshaling %#v: %s", object, err)
	}
	return o
}

func main() {
	nsqConfig.MaxInFlight = maxConcurrentHttpRequests
	initConsumers()
}

package main

import (
	"github.com/bitly/nsq/internal/app"
	"github.com/duskhacker/cqrsnu/cafe"
)

const (
	maxConcurrentHttpRequests = 5
)

func main() {
	cafe.SetLookupdHTTPAddrs = app.StringArray{}
	//	NsqConfig.MaxInFlight = maxConcurrentHttpRequests
	cafe.InitConsumers()
	// stopchan and all that
	cafe.StopAllConsumers()
}

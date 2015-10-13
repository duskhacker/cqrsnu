package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-nsq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pf = fmt.Printf

var _ = Describe("Main", func() {

	BeforeEach(func() {
		lookupdHTTPAddrs.Set("mnementh.dev:4161")
		nsqConfig.MaxInFlight = maxConcurrentHttpRequests
		nsqConfig.LookupdPollInterval = time.Millisecond * 100

		initHandlers()

	})

	Describe("Tab", func() {
		var (
			producer *nsq.Producer
			received = false
		)

		BeforeEach(func() {
			producer = newProducer(openTabTopic, openTabTopic+"Producer")
		})

		FIt("receives a TabOpened Event", func() {
			command := NewOpenTab(1, "Veronica")
			expected := NewTabOpened(command.Guid, 1, "Veronica")

			f := func(m *nsq.Message) error {
				to := TabOpened{}.FromJson(m.Body)
				if to.Guid.String() == expected.Guid.String() {
					received = true
				}
				return nil
			}

			_ = newConsumer(tabOpenedTopic, tabOpenedTopic+"TestConsumer", f)

			producer.Publish(openTabTopic, command.ToJson())

			Eventually(func() bool { return received }).Should(BeTrue(), "No TabOpened Event generated")
		})
	})
})

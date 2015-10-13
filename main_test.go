package main

import (
	"fmt"
	"log"

	"time"

	"github.com/bitly/go-nsq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pf = fmt.Printf

var _ = Describe("Main", func() {
	var (
		nsqdTCPAddr = "localhost:4150"
		consumer    *nsq.Consumer
		producer    *nsq.Producer
		err         error
		stuff       OpenTab
	)

	BeforeEach(func() {
		lookupdHTTPAddrs.Set("mnementh.dev:4161")
		cfg := nsq.NewConfig()
		cfg.UserAgent = channel
		cfg.MaxInFlight = maxConcurrentHttpRequests
		cfg.LookupdPollInterval = time.Millisecond * 250

		if consumer, err = nsq.NewConsumer("Tab", "TabConsumer", cfg); err != nil {
			log.Fatalf("NewConsumer: %s", err)
		}

		consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
			stuff = OpenTab{}.FromJson(msg.Body)
			return nil
		}))

		if err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs); err != nil {
			log.Fatalf("ConnectToNSQLookupds: %s", err)
		}

		cfg.UserAgent = fmt.Sprintf("%s_producer", "Tab")
		if producer, err = nsq.NewProducer(nsqdTCPAddr, cfg); err != nil {
			log.Fatalf("failed to create nsq.Producer - %s", err)
		}

		if err = producer.Ping(); err != nil {
			log.Fatalf("unable to ping nsqd: %s\n", err)
		}

	})

	Describe("Tab", func() {
		FIt("receives a TabOpened Event", func() {
			command := NewOpenTab(1, "Veronica")
			expected := NewTabOpened(command.Guid, 1, "Veronica")

			f := func() bool {
				for _, tab := range TabOpenedStore {
					if tab.Guid.String() == expected.Guid.String() {
						return true
					}
				}
				return false
			}

			producer.Publish(Topic, command.ToJson())

			Eventually(f).Should(BeTrue(), "No TabOpened Event generated")
		})
	})
})

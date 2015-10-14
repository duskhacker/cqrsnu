package main

import (
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"

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

		initConsumers()
	})

	Describe("Tab", func() {
		var (
			received = false
		)

		It("receives a TabOpened Event", func() {
			command := NewOpenTab(1, "Veronica")
			expected := NewTabOpened(command.ID, 1, "Veronica")

			f := func(m *nsq.Message) error {
				to := new(TabOpened).FromJson(m.Body)
				if to.ID.String() == expected.ID.String() {
					received = true
				}
				return nil
			}

			_ = newConsumer(tabOpened, tabOpened+"TestConsumer", f)

			Send(openTab, command)

			Eventually(func() bool { return received }).Should(BeTrue(), "No TabOpened Event generated")
		})
	})

	Describe("PlaceOrder", func() {
		var (
			received bool
			id       uuid.UUID
		)

		Describe("with no tab opened", func() {
			It("receives error if no tab opened", func() {
				command := NewPlaceOrder(uuid.NewRandom(), nil)

				f := func(m *nsq.Message) error {
					ex := new(CommandException).FromJson(m.Body)
					if ex.Type == "TabNotOpen" && ex.Message == "Cannot Place order without open Tab" {
						received = true
					}
					return nil
				}

				_ = newConsumer(exception, exception+"TestConsumer", f)

				Send(placeOrder, command)

				Eventually(func() bool { return received }).Should(BeTrue(), "PlaceOrder Exception not Raised")
			})
		})

		Describe("with tab opened", func() {
			var (
				foodOrderedReceived   bool
				drinksOrderedReceived bool

				drink OrderedItem
				food  OrderedItem
			)

			BeforeEach(func() {
				foodOrderedReceived = false
				drinksOrderedReceived = false

				drink = NewOrderedItem(1, "Patron", true, 5.00)
				food = NewOrderedItem(1, "Steak", false, 15.00)

				dof := func(m *nsq.Message) error {
					order := new(DrinksOrdered).FromJson(m.Body)
					if len(order.Items) > 0 {
						drinksOrderedReceived = true
					}
					return nil
				}

				_ = newConsumer(drinksOrdered, drinksOrdered+"TestConsumer", dof)

				fof := func(m *nsq.Message) error {
					order := new(FoodOrdered).FromJson(m.Body)
					if len(order.Items) > 0 {
						foodOrderedReceived = true
					}
					return nil
				}

				_ = newConsumer(foodOrdered, foodOrdered+"TestConsumer", fof)

				exf := func(m *nsq.Message) error {
					defer GinkgoRecover()
					ex := new(CommandException).FromJson(m.Body)
					Expect(ex).To(BeNil())
					return nil
				}

				_ = newConsumer(exception, exception+"TestConsumer", exf)

				command := NewOpenTab(1, "Veronica")
				id = command.ID

				Send(openTab, command)
			})

			It("orders drinks", func() {
				command := NewPlaceOrder(id, []OrderedItem{drink})

				Send(placeOrder, command)

				Eventually(func() bool { return drinksOrderedReceived }).Should(BeTrue(), "No DrinksOrdered event generated")
			})

			It("orders food", func() {
				command := NewPlaceOrder(id, []OrderedItem{food})

				Send(placeOrder, command)

				Eventually(func() bool { return foodOrderedReceived }).Should(BeTrue(), "No FoodOrdered event generated")
			})

			It("orders food and drink", func() {
				command := NewPlaceOrder(id, []OrderedItem{food, drink})

				Send(placeOrder, command)

				Eventually(func() bool { return foodOrderedReceived }).Should(BeTrue(), "No FoodOrdered event generated")
				Eventually(func() bool { return drinksOrderedReceived }).Should(BeTrue(), "No DrinksOrdered event generated")
			})
		})
	})
})

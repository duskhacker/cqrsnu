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
		It("opens a tab", func() {
			done := make(chan bool)
			command := NewOpenTab(1, "Veronica")
			expected := NewTabOpened(command.ID, 1, "Veronica")

			f := func(m *nsq.Message) error {
				defer GinkgoRecover()
				Expect(new(TabOpened).FromJson(m.Body)).To(Equal(expected))
				done <- true
				return nil
			}

			c := newConsumer(tabOpened, tabOpened+"TestConsumer", f)

			Send(openTab, command)

			Eventually(done).Should(Receive(BeTrue()), "No TabOpened Event generated")
			c.Stop()
		})
	})

	Describe("PlaceOrder", func() {
		Describe("with no tab opened", func() {
			done := make(chan bool)
			It("receives error", func() {
				command := NewPlaceOrder(nil, nil)
				expected := NewCommandException(nil, "TabNotOpen", "Cannot Place order without open Tab")

				f := func(m *nsq.Message) error {
					ex := new(CommandException).FromJson(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(expected))
					done <- true
					return nil
				}

				c := newConsumer(exception, exception+"TestConsumer", f)

				Send(placeOrder, command)

				Eventually(done).Should(Receive(BeTrue()), "PlaceOrder Exception not Raised")
				c.Stop()
			})
		})

		Describe("with tab opened", func() {
			var (
				c1, c2, c3 *nsq.Consumer
				id         uuid.UUID

				foodOrderedDone   = make(chan bool)
				drinksOrderedDone = make(chan bool)

				drink OrderedItem
				food  OrderedItem
			)

			BeforeEach(func() {
				drink = NewOrderedItem(1, "Patron", true, 5.00)
				food = NewOrderedItem(1, "Steak", false, 15.00)

				dof := func(m *nsq.Message) error {
					order := new(DrinksOrdered).FromJson(m.Body)
					if len(order.Items) > 0 {
						drinksOrderedDone <- true
					}
					return nil
				}

				c1 = newConsumer(drinksOrdered, drinksOrdered+"TestConsumer", dof)

				fof := func(m *nsq.Message) error {
					order := new(FoodOrdered).FromJson(m.Body)
					if len(order.Items) > 0 {
						foodOrderedDone <- true
					}
					return nil
				}

				c2 = newConsumer(foodOrdered, foodOrdered+"TestConsumer", fof)

				exf := func(m *nsq.Message) error {
					defer GinkgoRecover()
					ex := new(CommandException).FromJson(m.Body)
					Expect(ex).To(BeNil())
					return nil
				}

				c3 = newConsumer(exception, exception+"TestConsumer", exf)

				command := NewOpenTab(1, "Veronica")
				id = command.ID

				Send(openTab, command)
			})

			AfterEach(func() {
				c1.Stop()
				c2.Stop()
				c3.Stop()
			})

			It("orders drinks", func() {
				Send(placeOrder, NewPlaceOrder(id, []OrderedItem{drink}))

				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "No DrinksOrdered event generated")
			})

			It("orders food", func() {
				command := NewPlaceOrder(id, []OrderedItem{food})

				Send(placeOrder, command)

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "No FoodOrdered event generated")
			})

			It("orders food and drink", func() {
				command := NewPlaceOrder(id, []OrderedItem{food, drink})

				Send(placeOrder, command)

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "No FoodOrdered event generated")
				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "No DrinksOrdered event generated")
			})
		})
	})
})

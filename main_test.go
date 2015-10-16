package main

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"

	"github.com/bitly/go-nsq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pf = fmt.Printf

var _ = Describe("Main", func() {
	BeforeEach(func() {
		Tabs = NewTabs()
		nsqConfig.MaxInFlight = maxConcurrentHttpRequests
		connectToNSQD = true
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
			It("receives error", func() {
				done := make(chan bool)
				command := NewPlaceOrder(nil, nil)

				f := func(m *nsq.Message) error {
					ex := new(CommandException).FromJson(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(TabNotOpenException))
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

	Describe("Serving Drinks", func() {
		var (
			id uuid.UUID

			drinks []OrderedItem
		)

		BeforeEach(func() {
			drinks = append(drinks, NewOrderedItem(1, "Patron", true, 5.00))
			drinks = append(drinks, NewOrderedItem(2, "Scotch", true, 7.00))

			command := NewOpenTab(1, "Veronica")
			id = command.ID

			Send(openTab, command)
		})

		Describe("with 1 drink ordered", func() {
			BeforeEach(func() {
				Send(placeOrder, NewPlaceOrder(id, []OrderedItem{drinks[0]}))
			})

			It("generates exception if second drink is marked served", func() {
				done := make(chan bool)

				f := func(m *nsq.Message) error {
					ex := new(CommandException).FromJson(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(DrinksNotOutstanding))
					done <- true
					return nil
				}
				c := newConsumer(exception, exception+"TestConsumer", f)

				Send(markDrinksServed, NewMarkDrinksServed(id, []int{drinks[1].MenuNumber}))

				Eventually(done).Should(Receive(BeTrue()), "drinksNotOutstanding Exception not Raised")
				c.Stop()
			})
		})

		Describe("with drinks ordered", func() {
			var (
				menuNumbers      []int
				drinksServedDone chan bool
			)

			BeforeEach(func() {
				drinksServedDone = make(chan bool)
				for _, drink := range drinks {
					menuNumbers = append(menuNumbers, drink.MenuNumber)
				}

				Send(placeOrder, NewPlaceOrder(id, drinks))
			})

			It("marks drinks served", func() {
				dsf := func(m *nsq.Message) error {
					defer GinkgoRecover()
					evt := new(DrinksServed).FromJson(m.Body)
					Expect(evt.MenuNumbers).To(Equal(menuNumbers))
					drinksServedDone <- true
					return nil
				}

				c1 := newConsumer(drinksServed, drinksServed+"TestConsumer", dsf)

				Send(markDrinksServed, NewMarkDrinksServed(id, menuNumbers))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "No DrinksServed event generated")
				c1.Stop()
			})

			It("does not allow drinks to be served twice", func() {
				gotException := make(chan bool)
				exf := func(m *nsq.Message) error {
					defer GinkgoRecover()
					ex := new(CommandException).FromJson(m.Body)
					Expect(ex).To(Equal(DrinksNotOutstanding))
					gotException <- true
					return nil
				}
				c1 := newConsumer(exception, exception+"TestExceptionConsumer", exf)

				dsf := func(m *nsq.Message) error {
					defer GinkgoRecover()
					evt := new(DrinksServed).FromJson(m.Body)
					Expect(evt.MenuNumbers).To(Equal(menuNumbers))
					drinksServedDone <- true
					return nil
				}
				c2 := newConsumer(drinksServed, drinksServed+"TestConsumer", dsf)

				Send(markDrinksServed, NewMarkDrinksServed(id, menuNumbers))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "No DrinksServed event generated")

				Send(markDrinksServed, NewMarkDrinksServed(id, menuNumbers))
				Eventually(gotException).Should(Receive(BeTrue()), "No Exception raised")

				c1.Stop()
				c2.Stop()
			})
		})

	})
})

func listenForUnexpectedException() {
	f := func(m *nsq.Message) error {
		pf("EXCEPTION: %#v\n", new(CommandException).FromJson(m.Body))
		return nil
	}
	_ = newConsumer(exception, exception+"UnexpectedExceptionConsumer", f)
}

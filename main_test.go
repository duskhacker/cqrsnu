package main

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"

	"github.com/bitly/go-nsq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pf = fmt.Printf

var testConsumers []*nsq.Consumer

var _ = Describe("Main", func() {

	var (
		openTabCmd OpenTab
		tabID      uuid.UUID
		drinks     []OrderedItem
		food       []OrderedItem
	)

	BeforeEach(func() {
		Tabs = NewTabs()
		openTabCmd = NewOpenTab(1, "Kinessa")
		tabID = openTabCmd.ID

		drinks = []OrderedItem{}
		drinks = append(drinks, NewOrderedItem(1, "Patron", true, 5.00))
		drinks = append(drinks, NewOrderedItem(2, "Scotch", true, 3.50))

		food = []OrderedItem{}
		food = append(food, NewOrderedItem(1, "Steak", false, 15.00))
		food = append(food, NewOrderedItem(2, "Burger", false, 8.00))
	})

	AfterEach(func() {
		stopallTestConsumers()
	})

	Describe("Tab", func() {
		It("opens a tab", func() {
			done := make(chan bool)

			newTestConsumer(tabOpened, tabOpened+"TestConsumer",
				func(m *nsq.Message) error {
					defer GinkgoRecover()
					Expect(new(TabOpened).FromJSON(m.Body)).To(Equal(NewTabOpened(tabID, 1, "Kinessa")))
					done <- true
					return nil
				})

			Send(openTab, openTabCmd)

			Eventually(done).Should(Receive(BeTrue()), "No TabOpened received")
		})
	})

	Describe("Ordering", func() {
		Describe("with no tab opened", func() {
			It("receives error", func() {
				done := make(chan bool)
				command := NewPlaceOrder(nil, nil)

				newTestConsumer(exception, exception+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						Expect(new(Exception).FromJSON(m.Body)).To(Equal(TabNotOpenException))
						done <- true
						return nil
					})

				Send(placeOrder, command)

				Eventually(done).Should(Receive(BeTrue()), "TabNotOpenException Exception not Raised")
			})
		})

		Describe("with tab opened", func() {
			var (
				foodOrderedDone   = make(chan bool)
				drinksOrderedDone = make(chan bool)
			)

			BeforeEach(func() {

				newTestConsumer(drinksOrdered, drinksOrdered+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(DrinksOrdered).FromJSON(m.Body)
						if len(order.Items) > 0 {
							drinksOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(foodOrdered, foodOrdered+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(FoodOrdered).FromJSON(m.Body)
						if len(order.Items) > 0 {
							foodOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(exception, exception+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(Exception).FromJSON(m.Body)
						Expect(ex).To(BeNil())
						return nil
					})

				Send(openTab, openTabCmd)
			})

			It("drinks", func() {

				Send(placeOrder, NewPlaceOrder(tabID, drinks))

				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})

			It("food", func() {
				Send(placeOrder, NewPlaceOrder(tabID, food))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
			})

			It("food and drink", func() {
				Send(placeOrder, NewPlaceOrder(tabID, append(food, drinks...)))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})
		})
	})

	Describe("Serving Drinks", func() {
		BeforeEach(func() {
			Send(openTab, openTabCmd)
		})

		Describe("with 1 drink ordered", func() {
			BeforeEach(func() {
				Send(placeOrder, NewPlaceOrder(tabID, drinks[:1]))
			})

			It("generates exception if second drink is marked served", func() {
				done := make(chan bool)

				newTestConsumer(exception, exception+"TestConsumer", func(m *nsq.Message) error {
					ex := new(Exception).FromJSON(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(DrinksNotOutstanding))
					done <- true
					return nil
				})

				Send(markDrinksServed, NewMarkDrinksServed(tabID, drinks[1:2]))

				Eventually(done).Should(Receive(BeTrue()), "DrinksNotOutstanding Exception not Raised")
			})
		})

		Describe("with drinks ordered", func() {
			var (
				drinksServedDone chan bool
			)

			BeforeEach(func() {
				drinksServedDone = make(chan bool)

				Send(placeOrder, NewPlaceOrder(tabID, drinks))
			})

			It("marks drinks served", func() {

				newTestConsumer(drinksServed, drinksServed+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(DrinksServed).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(drinks))
						drinksServedDone <- true
						return nil
					})

				Send(markDrinksServed, NewMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")
			})

			It("does not allow drinks to be served twice", func() {
				rcvdException := make(chan bool)
				newTestConsumer(exception, exception+"TestExceptionConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(Exception).FromJSON(m.Body)
						Expect(ex).To(Equal(DrinksNotOutstanding))
						rcvdException <- true
						return nil
					})

				newTestConsumer(drinksServed, drinksServed+"TestConsumer", func(m *nsq.Message) error {
					defer GinkgoRecover()
					evt := new(DrinksServed).FromJSON(m.Body)
					Expect(evt.Items).To(Equal(drinks))
					drinksServedDone <- true
					return nil
				})

				Send(markDrinksServed, NewMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")

				Send(markDrinksServed, NewMarkDrinksServed(tabID, drinks))
				Eventually(rcvdException).Should(Receive(BeTrue()), "DrinksNotOutstanding exception not received")
			})
		})

	})

	Describe("Food", func() {

		BeforeEach(func() {
			Send(openTab, openTabCmd)
			Send(placeOrder, NewPlaceOrder(tabID, food))
		})

		Describe("prepare", func() {
			It("marks food prepared", func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(foodPrepared, foodPrepared+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(FoodPrepared).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodPrepared <- true
						return nil
					})

				Send(markFoodPrepared, NewMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")
			})
		})

		Describe("serve", func() {
			BeforeEach(func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(foodPrepared, foodPrepared+"TestConsumer",
					func(m *nsq.Message) error {
						rcvdFoodPrepared <- true
						return nil
					})

				Send(markFoodPrepared, NewMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")

			})

			It("marks food served", func() {
				listenForUnexpectedException()
				rcvdFoodServed := make(chan bool)
				newTestConsumer(foodServed, foodServed+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(FoodServed).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodServed <- true
						return nil
					})

				Send(markFoodServed, NewMarkFoodServed(tabID, food))
				Eventually(rcvdFoodServed).Should(Receive(BeTrue()), "FoodServed not received")
			})
		})
	})

	Describe("Closing Tab", func() {
		BeforeEach(func() {
			Send(openTab, openTabCmd)
			Send(placeOrder, NewPlaceOrder(tabID, drinks))
			Send(markDrinksServed, NewMarkDrinksServed(tabID, drinks))

			rcvdDrinksServed := make(chan bool)
			newTestConsumer(drinksServed, drinksServed+"TestConsumer",
				func(msg *nsq.Message) error {
					rcvdDrinksServed <- true
					return nil
				})

			Eventually(rcvdDrinksServed).Should(Receive(BeTrue()))
		})

		Describe("with tip", func() {
			It("closes tab", func() {
				tabClosedReceived := make(chan bool)
				newTestConsumer(tabClosed, tabClosed+"TestConsumer",
					func(msg *nsq.Message) error {
						evt := new(TabClosed).FromJSON(msg.Body)
						defer GinkgoRecover()
						Expect(evt.AmountPaid).To(Equal(8.50 + 0.50))
						Expect(evt.OrderValue).To(Equal(8.50))
						Expect(evt.TipValue).To(Equal(0.5))
						tabClosedReceived <- true
						return nil
					})

				Send(closeTab, NewCloseTab(tabID, 8.50+0.50))

				Eventually(tabClosedReceived).Should(Receive(BeTrue()), "TabClosed not received")
			})
		})
	})
})

func newTestConsumer(topic, channel string, f func(*nsq.Message) error) {
	testConsumers = append(testConsumers, newConsumer(topic, channel, f))
}

func stopallTestConsumers() {
	for _, consumer := range testConsumers {
		consumer.Stop()
	}
}

func listenForUnexpectedException() {
	f := func(m *nsq.Message) error {
		pf("EXCEPTION: %#v\n", new(Exception).FromJSON(m.Body))
		return nil
	}
	newTestConsumer(exception, exception+"UnexpectedExceptionConsumer", f)
}

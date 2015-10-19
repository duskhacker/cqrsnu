package cafe

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
		openTabCmd openTab
		tabID      uuid.UUID
		drinks     []OrderedItem
		food       []OrderedItem
	)

	BeforeEach(func() {
		Tabs = NewTabs()
		openTabCmd = newOpenTab(1, "Kinessa")
		tabID = openTabCmd.ID

		drinks = []OrderedItem{}
		drinks = append(drinks, NewOrderedItem(1, "Patron", true, 5.00))
		drinks = append(drinks, NewOrderedItem(2, "Scotch", true, 3.50))

		food = []OrderedItem{}
		food = append(food, NewOrderedItem(1, "Steak", false, 15.00))
		food = append(food, NewOrderedItem(2, "Burger", false, 8.00))
	})

	AfterEach(func() {
		stopAllTestConsumers()
	})

	Describe("Tab", func() {
		It("opens a tab", func() {
			done := make(chan bool)

			newTestConsumer(tabOpenedTopic, tabOpenedTopic+"TestConsumer",
				func(m *nsq.Message) error {
					defer GinkgoRecover()
					Expect(new(tabOpened).fromJSON(m.Body)).To(Equal(newTabOpened(tabID, 1, "Kinessa")))
					done <- true
					return nil
				})

			Send(openTabTopic, openTabCmd)

			Eventually(done).Should(Receive(BeTrue()), "No TabOpened received")
		})
	})

	Describe("Ordering", func() {
		Describe("with no tab opened", func() {
			It("receives error", func() {
				done := make(chan bool)
				command := newPlaceOrder(nil, nil)

				newTestConsumer(exceptionTopic, exceptionTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						Expect(new(exception).fromJSON(m.Body)).To(Equal(TabNotOpenException))
						done <- true
						return nil
					})

				Send(placeOrderTopic, command)

				Eventually(done).Should(Receive(BeTrue()), "TabNotOpenException Exception not Raised")
			})
		})

		Describe("with tab opened", func() {
			var (
				foodOrderedDone   = make(chan bool)
				drinksOrderedDone = make(chan bool)
			)

			BeforeEach(func() {

				newTestConsumer(drinksOrderedTopic, drinksOrderedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(drinksOrdered).fromJSON(m.Body)
						if len(order.Items) > 0 {
							drinksOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(foodOrderedTopic, foodOrderedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(foodOrdered).fromJSON(m.Body)
						if len(order.Items) > 0 {
							foodOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(exceptionTopic, exceptionTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(exception).fromJSON(m.Body)
						Expect(ex).To(BeNil())
						return nil
					})

				Send(openTabTopic, openTabCmd)
			})

			It("drinks", func() {

				Send(placeOrderTopic, newPlaceOrder(tabID, drinks))

				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})

			It("food", func() {
				Send(placeOrderTopic, newPlaceOrder(tabID, food))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
			})

			It("food and drink", func() {
				Send(placeOrderTopic, newPlaceOrder(tabID, append(food, drinks...)))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})
		})
	})

	Describe("Serving Drinks", func() {
		BeforeEach(func() {
			Send(openTabTopic, openTabCmd)
		})

		Describe("with 1 drink ordered", func() {
			BeforeEach(func() {
				Send(placeOrderTopic, newPlaceOrder(tabID, drinks[:1]))
			})

			It("generates exception if second drink is marked served", func() {
				done := make(chan bool)

				newTestConsumer(exceptionTopic, exceptionTopic+"TestConsumer", func(m *nsq.Message) error {
					ex := new(exception).fromJSON(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(DrinksNotOutstanding))
					done <- true
					return nil
				})

				Send(markDrinksServedTopic, newMarkDrinksServed(tabID, drinks[1:2]))

				Eventually(done).Should(Receive(BeTrue()), "DrinksNotOutstanding Exception not Raised")
			})
		})

		Describe("with drinks ordered", func() {
			var (
				drinksServedDone chan bool
			)

			BeforeEach(func() {
				drinksServedDone = make(chan bool)

				Send(placeOrderTopic, newPlaceOrder(tabID, drinks))
			})

			It("marks drinks served", func() {

				newTestConsumer(drinksServedTopic, drinksServedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(drinksServed).fromJSON(m.Body)
						Expect(evt.Items).To(Equal(drinks))
						drinksServedDone <- true
						return nil
					})

				Send(markDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")
			})

			It("does not allow drinks to be served twice", func() {
				rcvdException := make(chan bool)
				newTestConsumer(exceptionTopic, exceptionTopic+"TestExceptionConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(exception).fromJSON(m.Body)
						Expect(ex).To(Equal(DrinksNotOutstanding))
						rcvdException <- true
						return nil
					})

				newTestConsumer(drinksServedTopic, drinksServedTopic+"TestConsumer", func(m *nsq.Message) error {
					defer GinkgoRecover()
					evt := new(drinksServed).fromJSON(m.Body)
					Expect(evt.Items).To(Equal(drinks))
					drinksServedDone <- true
					return nil
				})

				Send(markDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")

				Send(markDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(rcvdException).Should(Receive(BeTrue()), "DrinksNotOutstanding exception not received")
			})
		})

	})

	Describe("Food", func() {

		BeforeEach(func() {
			Send(openTabTopic, openTabCmd)
			Send(placeOrderTopic, newPlaceOrder(tabID, food))
		})

		Describe("prepare", func() {
			It("marks food prepared", func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(foodPreparedTopic, foodPreparedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(foodPrepared).fromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodPrepared <- true
						return nil
					})

				Send(markFoodPreparedTopic, newMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")
			})
		})

		Describe("serve", func() {
			BeforeEach(func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(foodPreparedTopic, foodPreparedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						rcvdFoodPrepared <- true
						return nil
					})

				Send(markFoodPreparedTopic, newMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")

			})

			It("marks food served", func() {
				listenForUnexpectedException()
				rcvdFoodServed := make(chan bool)
				newTestConsumer(foodServedTopic, foodServedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(foodServed).fromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodServed <- true
						return nil
					})

				Send(markFoodServedTopic, newMarkFoodServed(tabID, food))
				Eventually(rcvdFoodServed).Should(Receive(BeTrue()), "FoodServed not received")
			})
		})
	})

	Describe("Closing Tab", func() {
		BeforeEach(func() {
			Send(openTabTopic, openTabCmd)
			Send(placeOrderTopic, newPlaceOrder(tabID, append(food, drinks...)))
			Send(markDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
			Send(markFoodPreparedTopic, newMarkFoodPrepared(tabID, food))
			Send(markFoodServedTopic, newMarkFoodServed(tabID, food))

			rcvdDrinksServed := make(chan bool)
			newTestConsumer(drinksServedTopic, drinksServedTopic+"TestConsumer",
				func(msg *nsq.Message) error {
					rcvdDrinksServed <- true
					return nil
				})

			Eventually(rcvdDrinksServed).Should(Receive(BeTrue()))
		})

		Describe("with tip", func() {
			It("closes tab", func() {
				tabClosedReceived := make(chan bool)
				newTestConsumer(tabClosedTopic, tabClosedTopic+"TestConsumer",
					func(msg *nsq.Message) error {
						evt := new(tabClosed).fromJSON(msg.Body)
						defer GinkgoRecover()
						Expect(evt.AmountPaid).To(Equal(31.50 + 0.50))
						Expect(evt.OrderValue).To(Equal(31.50))
						Expect(evt.TipValue).To(Equal(0.5))
						tabClosedReceived <- true
						return nil
					})

				Send(closeTabTopic, newCloseTab(tabID, 31.50+0.50))

				Eventually(tabClosedReceived).Should(Receive(BeTrue()), "TabClosed not received")
			})
		})
	})
})

func newTestConsumer(topic, channel string, f func(*nsq.Message) error) {
	testConsumers = append(testConsumers, NewConsumer(topic, channel, f))
}

func stopAllTestConsumers() {
	for _, consumer := range testConsumers {
		consumer.Stop()
	}
}

func listenForUnexpectedException() {
	f := func(m *nsq.Message) error {
		pf("EXCEPTION: %#v\n", new(exception).fromJSON(m.Body))
		return nil
	}
	newTestConsumer(exceptionTopic, exceptionTopic+"UnexpectedExceptionConsumer", f)
}

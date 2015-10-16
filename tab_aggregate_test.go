package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TabAggegate", func() {
	Describe("DeleteOrderedItem", func() {

		BeforeEach(func() {
			Tabs = NewTabs()
		})

		It("deletes an item", func() {
			drink0 := NewOrderedItem(5, "drink1", true, 0.0)
			drink1 := NewOrderedItem(1, "drink2", true, 0.0)
			drink2 := NewOrderedItem(7, "drink3", true, 0.0)

			drinks := []OrderedItem{}
			drinks = append(drinks, drink0)
			drinks = append(drinks, drink1)
			drinks = append(drinks, drink2)

			tab := NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[1:2])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink0, drink2}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[0:1])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink1, drink2}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[2:])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink0, drink1}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks)
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{}))
		})
	})

})

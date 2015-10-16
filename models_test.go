package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("DeleteOrderedItem", func() {

		BeforeEach(func() {
			Tabs = NewTabs()
		})

		It("deletes an item", func() {
			drink1 := NewOrderedItem(5, "drink1", true, 0.0)
			drink2 := NewOrderedItem(1, "drink2", true, 0.0)
			drink3 := NewOrderedItem(7, "drink3", true, 0.0)

			drinks := []OrderedItem{}
			drinks = append(drinks, drink1)
			drinks = append(drinks, drink2)
			drinks = append(drinks, drink3)

			tab := NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks([]int{1})
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink1, drink3}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks([]int{5})
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink2, drink3}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks([]int{7})
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink1, drink2}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks([]int{1, 5, 7})
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{}))
		})
	})

})

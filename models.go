package main

import "code.google.com/p/go-uuid/uuid"

type Tab struct {
	ID                uuid.UUID
	TableNumber       int
	WaitStaff         string
	OutstandingDrinks []OrderedItem
	OutstandingFood   []OrderedItem
	Open              bool
	ServedItemsValue  float64
}

func (t Tab) AreDrinksOutstanding(drinks []int) bool {
	for _, drink := range drinks {
		for _, outstanding := range t.OutstandingDrinks {
			if drink == outstanding.MenuNumber {
				return true
			}
		}
	}
	return false
}

// -

type OrderedItem struct {
	MenuNumber  int
	Description string
	IsDrink     bool
	Price       float64
}

func NewOrderedItem(menuNumber int, description string, isDrink bool, price float64) OrderedItem {
	return OrderedItem{
		MenuNumber:  menuNumber,
		Description: description,
		IsDrink:     isDrink,
		Price:       price,
	}
}

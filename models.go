package main

import (
	"fmt"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

var (
	Tabs map[string]*Tab
)

func NewTabs() map[string]*Tab {
	return make(map[string]*Tab)
}

type Tab struct {
	ID                uuid.UUID
	TableNumber       int
	WaitStaff         string
	OutstandingDrinks []OrderedItem
	OutstandingFood   []OrderedItem
	Open              bool
	ServedItemsValue  float64
}

func NewTab(id uuid.UUID, tn int, ws string, od []OrderedItem, of []OrderedItem, open bool, siv float64) *Tab {
	mutex.Lock()
	defer mutex.Unlock()
	tab := &Tab{
		ID:                id,
		TableNumber:       tn,
		WaitStaff:         ws,
		OutstandingDrinks: od,
		OutstandingFood:   of,
		Open:              open,
		ServedItemsValue:  siv,
	}
	Tabs[id.String()] = tab
	return tab
}

func GetTab(id uuid.UUID) *Tab {
	tab, ok := Tabs[id.String()]
	if !ok {
		return nil
	}
	return tab
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

func indexOfOrderedItem(items []OrderedItem, menuNumber int) int {
	for i, e := range items {
		if e.MenuNumber == menuNumber {
			return i
		}
	}
	return -1
}

func (t *Tab) DeleteOutstandingDrinks(items []int) error {
	for _, item := range items {
		drinks, err := deleteOrderedItem(t.OutstandingDrinks, item)
		if err != nil {
			return err
		}
		t.OutstandingDrinks = drinks
	}
	return nil
}

func deleteOrderedItem(items []OrderedItem, item int) ([]OrderedItem, error) {
	idx := indexOfOrderedItem(items, item)
	if idx < 0 {
		return nil, fmt.Errorf("no item %#v in tab", item)
	}
	a := make([]OrderedItem, len(items))
	n := copy(a, items)
	if n <= 0 {
		log.Fatalf("error copying data for deleteOutstandingDrinks")
	}
	return append(a[:idx], a[idx+1:]...), nil
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

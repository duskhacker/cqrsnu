package main

import "github.com/bitly/go-nsq"

func OpenTabHandler(msg *nsq.Message) error {
	ot := OpenTab{}.FromJson(msg.Body)
	Tabs[ot.ID.String()] = &Tab{ID: ot.ID, TableNumber: ot.TableNumber, WaitStaff: ot.WaitStaff}
	Send(tabOpened, NewTabOpened(ot.ID, ot.TableNumber, ot.WaitStaff))
	return nil
}

func PlaceOrderHandler(msg *nsq.Message) error {
	order := new(PlaceOrder).FromJson(msg.Body)
	tab, ok := Tabs[order.ID.String()]
	if !ok {
		Send(exception, tabNotOpenException)
		return nil
	}

	var (
		foodItems  []OrderedItem
		drinkItems []OrderedItem
	)

	for _, item := range order.Items {
		if item.IsDrink {
			drinkItems = append(drinkItems, item)
		} else {
			foodItems = append(foodItems, item)
		}
	}

	if len(foodItems) > 0 {
		tab.OutstandingFood = append(tab.OutstandingFood, foodItems...)
		Send(foodOrdered, NewFoodOrdered(order.ID, foodItems))
	}

	if len(drinkItems) > 0 {
		tab.OutstandingDrinks = append(tab.OutstandingDrinks, drinkItems...)
		Send(drinksOrdered, NewDrinksOrdered(order.ID, drinkItems))
	}

	return nil
}

func MarkDrinksServedHandler(msg *nsq.Message) error {
	c := new(MarkDrinksServed).FromJson(msg.Body)
	tab, ok := Tabs[c.ID.String()]
	if !ok {
		Send(exception, tabNotOpenException)
		return nil
	}

	if !tab.AreDrinksOutstanding(c.MenuNumbers) {
		Send(exception, drinksNotOutstanding)
		return nil
	}

	Send(drinksServed, NewDrinksServed(c.ID, c.MenuNumbers))
	return nil
}

func DrinksServedHandler(msg *nsq.Message) error {
	c := new(DrinksServed).FromJson(msg.Body)
	tab, ok := Tabs[c.ID.String()]
	if !ok {
		Send(exception, tabNotOpenException)
		return nil
	}

	tab.DeleteOutstandingDrinks(c.MenuNumbers)

	return nil
}

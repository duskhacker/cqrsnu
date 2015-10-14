package main

import "github.com/bitly/go-nsq"

func OpenTabHandler(msg *nsq.Message) error {
	ot := OpenTab{}.FromJson(msg.Body)
	Tabs[ot.ID.String()] = Tab{TableNumber: ot.TableNumber, WaitStaff: ot.WaitStaff}
	Send(tabOpened, NewTabOpened(ot.ID, ot.TableNumber, ot.WaitStaff))
	return nil
}

func PlaceOrderHandler(msg *nsq.Message) error {
	order := new(PlaceOrder).FromJson(msg.Body)
	tab, ok := Tabs[order.ID.String()]
	if !ok {
		Send(exception, NewCommandException(nil, "TabNotOpen", "Cannot Place order without open Tab"))
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
		Send(foodOrdered, NewFoodOrdered(order.ID, foodItems))
	}

	if len(drinkItems) > 0 {
		Send(drinksOrdered, NewDrinksOrdered(order.ID, drinkItems))
	}

	tab.Items = append(tab.Items, order.Items...)

	return nil
}

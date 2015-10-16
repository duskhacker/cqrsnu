package main

import "github.com/bitly/go-nsq"

func OpenTabHandler(msg *nsq.Message) error {
	cmd := new(OpenTab).FromJSON(msg.Body)
	NewTab(cmd.ID, cmd.TableNumber, cmd.WaitStaff, nil, nil, true, 0)
	Send(tabOpened, NewTabOpened(cmd.ID, cmd.TableNumber, cmd.WaitStaff))
	return nil
}

func PlaceOrderHandler(msg *nsq.Message) error {
	cmd := new(PlaceOrder).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	var (
		foodItems  []OrderedItem
		drinkItems []OrderedItem
	)

	for _, item := range cmd.Items {
		if item.IsDrink {
			drinkItems = append(drinkItems, item)
		} else {
			foodItems = append(foodItems, item)
		}
	}

	if len(foodItems) > 0 {
		tab.OutstandingFoods = append(tab.OutstandingFoods, foodItems...)
		Send(foodOrdered, NewFoodOrdered(cmd.ID, foodItems))
	}

	if len(drinkItems) > 0 {
		tab.OutstandingDrinks = append(tab.OutstandingDrinks, drinkItems...)
		Send(drinksOrdered, NewDrinksOrdered(cmd.ID, drinkItems))
	}

	return nil
}

func MarkDrinksServedHandler(msg *nsq.Message) error {
	cmd := new(MarkDrinksServed).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	if !tab.AreDrinksOutstanding(cmd.Items) {
		Send(exception, DrinksNotOutstanding)
		return nil
	}

	Send(drinksServed, NewDrinksServed(cmd.ID, cmd.Items))
	return nil
}

func MarkFoodServedHandler(msg *nsq.Message) error {
	cmd := new(MarkFoodServed).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	if !tab.AreFoodsOutstanding(cmd.Items) {
		Send(exception, FoodsNotOutstanding)
		return nil
	}

	Send(foodServed, NewFoodServed(cmd.ID, cmd.Items))
	return nil
}

func MarkFoodPreparedHandler(msg *nsq.Message) error {
	cmd := new(MarkFoodPrepared).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	if !tab.AreFoodsOutstanding(cmd.Items) {
		Send(exception, FoodsNotOutstanding)
		return nil
	}

	Send(foodPrepared, NewFoodPrepared(cmd.ID, cmd.Items))
	return nil
}

func DrinksServedHandler(msg *nsq.Message) error {
	evt := new(DrinksServed).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(evt.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	tab.AddServedItemsValue(evt.Items)
	tab.DeleteOutstandingDrinks(evt.Items)

	return nil
}

func FoodServedHandler(msg *nsq.Message) error {
	evt := new(FoodServed).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(evt.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	tab.AddServedItemsValue(evt.Items)
	tab.DeleteOutstandingFoods(evt.Items)

	return nil
}

func CloseTabHandler(msg *nsq.Message) error {
	cmd := new(CloseTab).FromJSON(msg.Body)
	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exception, TabNotOpenException)
		return nil
	}

	tipValue := cmd.AmountPaid - tab.ServedItemsValue

	Send(tabClosed, NewTabClosed(cmd.ID, cmd.AmountPaid, tab.ServedItemsValue, tipValue))
	return nil
}

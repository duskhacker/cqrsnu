package cafe

import "github.com/bitly/go-nsq"

func OpenTabHandler(msg *nsq.Message) error {
	cmd := new(openTab).fromJSON(msg.Body)
	NewTab(cmd.ID, cmd.TableNumber, cmd.WaitStaff, nil, nil, true, 0)
	Send(tabOpenedTopic, newTabOpened(cmd.ID, cmd.TableNumber, cmd.WaitStaff))
	return nil
}

func PlaceOrderHandler(msg *nsq.Message) error {
	cmd := new(placeOrder).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
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
		Send(foodOrderedTopic, newFoodOrdered(cmd.ID, foodItems))
	}

	if len(drinkItems) > 0 {
		tab.OutstandingDrinks = append(tab.OutstandingDrinks, drinkItems...)
		Send(drinksOrderedTopic, newDrinksOrdered(cmd.ID, drinkItems))
	}

	return nil
}

func MarkDrinksServedHandler(msg *nsq.Message) error {
	cmd := new(markDrinksServed).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	if !tab.AreDrinksOutstanding(cmd.Items) {
		Send(exceptionTopic, DrinksNotOutstanding)
		return nil
	}

	Send(drinksServedTopic, newDrinksServed(cmd.ID, cmd.Items))
	return nil
}

func MarkFoodServedHandler(msg *nsq.Message) error {
	cmd := new(markFoodServed).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	if !tab.AreFoodsOutstanding(cmd.Items) {
		Send(exceptionTopic, FoodsNotOutstanding)
		return nil
	}

	Send(foodServedTopic, newFoodServed(cmd.ID, cmd.Items))
	return nil
}

func MarkFoodPreparedHandler(msg *nsq.Message) error {
	cmd := new(markFoodPrepared).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	if !tab.AreFoodsOutstanding(cmd.Items) {
		Send(exceptionTopic, FoodsNotOutstanding)
		return nil
	}

	Send(foodPreparedTopic, newFoodPrepared(cmd.ID, cmd.Items))
	return nil
}

func DrinksServedHandler(msg *nsq.Message) error {
	evt := new(drinksServed).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(evt.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	tab.AddServedItemsValue(evt.Items)
	tab.DeleteOutstandingDrinks(evt.Items)

	return nil
}

func FoodServedHandler(msg *nsq.Message) error {
	evt := new(foodServed).fromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(evt.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	tab.AddServedItemsValue(evt.Items)
	tab.DeleteOutstandingFoods(evt.Items)

	return nil
}

func CloseTabHandler(msg *nsq.Message) error {
	cmd := new(closeTab).fromJSON(msg.Body)
	mutex.Lock()
	defer mutex.Unlock()

	tab := GetTab(cmd.ID)
	if tab == nil {
		Send(exceptionTopic, TabNotOpenException)
		return nil
	}

	tipValue := cmd.AmountPaid - tab.ServedItemsValue

	Send(tabClosedTopic, newTabClosed(cmd.ID, cmd.AmountPaid, tab.ServedItemsValue, tipValue))
	return nil
}

package chef_todo_list

import (
	"fmt"
	"log"

	"github.com/bitly/go-nsq"
	"github.com/duskhacker/cqrsnu/cafe"
)

var (
	consumers []*nsq.Consumer
)

func InitConsumers() {
	consumer := cafe.NewConsumer(cafe.FoodOrderedTopic, cafe.FoodOrderedTopic+"Todo", FoodOrderedHandler)
	consumers = append(consumers, consumer)
	consumer = cafe.NewConsumer(cafe.FoodPreparedTopic, cafe.FoodPreparedTopic+"Todo", FoodPreparedHandler)
	consumers = append(consumers, consumer)
}

func FoodOrderedHandler(msg *nsq.Message) error {
	evt := new(cafe.FoodOrdered).FromJSON(msg.Body)

	group := newTodoListGroup(evt.ID, evt.Items)

	mutex.Lock()
	defer mutex.Unlock()

	chefTodoList = append(chefTodoList, group)

	return nil
}

func FoodPreparedHandler(msg *nsq.Message) error {
	evt := new(cafe.FoodPrepared).FromJSON(msg.Body)

	mutex.Lock()
	defer mutex.Unlock()

	list := getTodoListGroup(evt.ID)
	if list == nil {
		return fmt.Errorf("error finding todolist group for %s\n", evt.ID)

	}

	for _, item := range evt.Items {
		items, err := deleteTodoListItem(list.Items, item)
		if err != nil {
			return err
		}
		list.Items = items
	}

	return nil
}

func deleteTodoListItem(items []todoListItem, item cafe.OrderedItem) ([]todoListItem, error) {
	idx := indexOfTodoListItem(items, item)
	if idx < 0 {
		return nil, fmt.Errorf("no item %#v in tab", item)
	}
	a := make([]todoListItem, len(items))
	n := copy(a, items)
	if n <= 0 {
		log.Fatalf("error copying data for deleteOutstandingDrinks")
	}
	return append(a[:idx], a[idx+1:]...), nil
}

func indexOfTodoListItem(items []todoListItem, item cafe.OrderedItem) int {
	for i, e := range items {
		if e.MenuNumber == item.MenuNumber {
			return i
		}
	}
	return -1
}

func StopAllConsumers() {
	for _, consumer := range consumers {
		consumer.Stop()
	}
}

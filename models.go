package main

import "code.google.com/p/go-uuid/uuid"

type Tab struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
	Items       []OrderedItem
}

type OrderedItem struct {
	MenuNumber  int
	Description string
	IsDrink     bool
	Price       float64
}

//func (oi OrderedItem) FromJson(data []byte) OrderedItem {
//	var err error
//	err = json.Unmarshal(data, &oi)
//	if err != nil {
//		log.Fatalf("json.Unmarshal: %s\n'", err)
//	}
//	return oi
//}

func NewOrderedItem(menuNumber int, description string, isDrink bool, price float64) OrderedItem {
	return OrderedItem{
		MenuNumber:  menuNumber,
		Description: description,
		IsDrink:     isDrink,
		Price:       price,
	}
}

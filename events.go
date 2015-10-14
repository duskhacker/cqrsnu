package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

type TabOpened struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func (t TabOpened) FromJson(data []byte) TabOpened {
	var err error
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return t
}

func NewTabOpened(guid uuid.UUID, tableNumber int, waitStaff string) TabOpened {
	return TabOpened{
		ID:          guid,
		TableNumber: tableNumber,
		WaitStaff:   waitStaff,
	}
}

// --

type DrinksOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (do DrinksOrdered) FromJson(data []byte) DrinksOrdered {
	var err error
	err = json.Unmarshal(data, &do)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return do
}

func NewDrinksOrdered(id uuid.UUID, items []OrderedItem) DrinksOrdered {
	return DrinksOrdered{
		ID:    id,
		Items: items,
	}
}

// --

type FoodOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fo FoodOrdered) FromJson(data []byte) FoodOrdered {
	var err error
	err = json.Unmarshal(data, &fo)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fo
}

func NewFoodOrdered(id uuid.UUID, items []OrderedItem) FoodOrdered {
	return FoodOrdered{
		ID:    id,
		Items: items,
	}
}

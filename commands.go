package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

const (
	openTab    = "OpenTab"
	placeOrder = "PlaceOrder"
)

type OpenTab struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func NewOpenTab(tableNumber int, waiter string) OpenTab {
	return OpenTab{
		ID:          uuid.NewRandom(),
		TableNumber: tableNumber,
		WaitStaff:   waiter,
	}
}

func (o OpenTab) FromJson(data []byte) OpenTab {
	var err error
	err = json.Unmarshal(data, &o)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return o
}

// --

type PlaceOrder struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (po PlaceOrder) FromJson(data []byte) PlaceOrder {
	var err error
	err = json.Unmarshal(data, &po)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return po
}

func NewPlaceOrder(id uuid.UUID, items []OrderedItem) PlaceOrder {
	return PlaceOrder{
		ID:    id,
		Items: items,
	}
}

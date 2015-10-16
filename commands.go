package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

const (
	openTab          = "OpenTab"
	placeOrder       = "PlaceOrder"
	markDrinksServed = "MarkDrinksServed"
	closeTab         = "CloseTab"
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

func (o OpenTab) FromJSON(data []byte) OpenTab {
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

func (po PlaceOrder) FromJSON(data []byte) PlaceOrder {
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

// --

type MarkDrinksServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mds MarkDrinksServed) FromJSON(data []byte) MarkDrinksServed {
	var err error
	err = json.Unmarshal(data, &mds)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mds
}

func NewMarkDrinksServed(id uuid.UUID, items []OrderedItem) MarkDrinksServed {
	return MarkDrinksServed{
		ID:    id,
		Items: items,
	}
}

// --

type CloseTab struct {
	ID         uuid.UUID
	AmountPaid float64
}

func (ct CloseTab) FromJSON(data []byte) CloseTab {
	var err error
	err = json.Unmarshal(data, &ct)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return ct
}

func NewCloseTab(id uuid.UUID, amountPaid float64) CloseTab {
	return CloseTab{
		ID:         id,
		AmountPaid: amountPaid,
	}
}

package cafe

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

const (
	tabOpenedTopic     = "TabOpened"
	foodOrderedTopic   = "FoodOrdered"
	drinksOrderedTopic = "DrinksOrdered"
	drinksServedTopic  = "DrinksServed"
	foodPreparedTopic  = "FoodPrepared"
	foodServedTopic    = "FoodServed"
	tabClosedTopic     = "TabClosed"
	exceptionTopic     = "Exception"
)

type tabOpened struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func (t tabOpened) fromJSON(data []byte) tabOpened {
	var err error
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return t
}

func newTabOpened(guid uuid.UUID, tableNumber int, waitStaff string) tabOpened {
	return tabOpened{
		ID:          guid,
		TableNumber: tableNumber,
		WaitStaff:   waitStaff,
	}
}

// --

type drinksOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (do drinksOrdered) fromJSON(data []byte) drinksOrdered {
	var err error
	err = json.Unmarshal(data, &do)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return do
}

func newDrinksOrdered(id uuid.UUID, items []OrderedItem) drinksOrdered {
	return drinksOrdered{
		ID:    id,
		Items: items,
	}
}

// --

type foodOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fo foodOrdered) fromJSON(data []byte) foodOrdered {
	var err error
	err = json.Unmarshal(data, &fo)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fo
}

func newFoodOrdered(id uuid.UUID, items []OrderedItem) foodOrdered {
	return foodOrdered{
		ID:    id,
		Items: items,
	}
}

// --

type drinksServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (ds drinksServed) fromJSON(data []byte) drinksServed {
	var err error
	err = json.Unmarshal(data, &ds)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return ds
}

func newDrinksServed(id uuid.UUID, items []OrderedItem) drinksServed {
	return drinksServed{
		ID:    id,
		Items: items,
	}
}

// --

type foodPrepared struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fp foodPrepared) fromJSON(data []byte) foodPrepared {
	var err error
	err = json.Unmarshal(data, &fp)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fp
}

func newFoodPrepared(id uuid.UUID, items []OrderedItem) foodPrepared {
	return foodPrepared{
		ID:    id,
		Items: items,
	}
}

// --

type foodServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fs foodServed) fromJSON(data []byte) foodServed {
	var err error
	err = json.Unmarshal(data, &fs)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fs
}

func newFoodServed(id uuid.UUID, items []OrderedItem) foodServed {
	return foodServed{
		ID:    id,
		Items: items,
	}
}

// --

type tabClosed struct {
	ID         uuid.UUID
	AmountPaid float64
	OrderValue float64
	TipValue   float64
}

func (tc tabClosed) fromJSON(data []byte) tabClosed {
	var err error
	err = json.Unmarshal(data, &tc)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return tc
}

func newTabClosed(id uuid.UUID, amountPaid, orderValue, tipValue float64) tabClosed {
	return tabClosed{
		ID:         id,
		AmountPaid: amountPaid,
		OrderValue: orderValue,
		TipValue:   tipValue,
	}
}

// --

type exception struct {
	Type    string
	Message string
}

func (e exception) fromJSON(data []byte) exception {
	var err error
	err = json.Unmarshal(data, &e)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return e
}

func newException(t string, msg string) exception {
	return exception{Type: t, Message: msg}
}

func (c exception) Error() string {
	return c.Type + ":" + c.Message
}

var TabNotOpenException = newException("TabNotOpen", "Cannot Place order without open Tab")
var DrinksNotOutstanding = newException("DrinksNotOutstanding", "Cannot serve unordered drinks")
var FoodsNotOutstanding = newException("FoodsNotOutstanding", "Cannot prepare unordered food")

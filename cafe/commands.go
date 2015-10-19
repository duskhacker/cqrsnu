package cafe

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

const (
	openTabTopic          = "OpenTab"
	placeOrderTopic       = "PlaceOrder"
	markDrinksServedTopic = "MarkDrinksServed"
	markFoodPreparedTopic = "MarkFoodPrepared"
	markFoodServedTopic   = "MarkFoodServed"
	closeTabTopic         = "CloseTab"
)

type openTab struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func newOpenTab(tableNumber int, waiter string) openTab {
	return openTab{
		ID:          uuid.NewRandom(),
		TableNumber: tableNumber,
		WaitStaff:   waiter,
	}
}

func (o openTab) fromJSON(data []byte) openTab {
	var err error
	err = json.Unmarshal(data, &o)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return o
}

// --

type placeOrder struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (po placeOrder) fromJSON(data []byte) placeOrder {
	var err error
	err = json.Unmarshal(data, &po)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return po
}

func newPlaceOrder(id uuid.UUID, items []OrderedItem) placeOrder {
	return placeOrder{
		ID:    id,
		Items: items,
	}
}

// --

type markDrinksServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mds markDrinksServed) fromJSON(data []byte) markDrinksServed {
	var err error
	err = json.Unmarshal(data, &mds)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mds
}

func newMarkDrinksServed(id uuid.UUID, items []OrderedItem) markDrinksServed {
	return markDrinksServed{
		ID:    id,
		Items: items,
	}
}

// --

type markFoodPrepared struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mfp markFoodPrepared) fromJSON(data []byte) markFoodPrepared {
	var err error
	err = json.Unmarshal(data, &mfp)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mfp
}

func newMarkFoodPrepared(id uuid.UUID, items []OrderedItem) markFoodPrepared {
	return markFoodPrepared{
		ID:    id,
		Items: items,
	}
}

// --

type markFoodServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mfs markFoodServed) fromJSON(data []byte) markFoodServed {
	var err error
	err = json.Unmarshal(data, &mfs)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mfs
}

func newMarkFoodServed(id uuid.UUID, items []OrderedItem) markFoodServed {
	return markFoodServed{
		ID:    id,
		Items: items,
	}
}

// --

type closeTab struct {
	ID         uuid.UUID
	AmountPaid float64
}

func (ct closeTab) fromJSON(data []byte) closeTab {
	var err error
	err = json.Unmarshal(data, &ct)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return ct
}

func newCloseTab(id uuid.UUID, amountPaid float64) closeTab {
	return closeTab{
		ID:         id,
		AmountPaid: amountPaid,
	}
}

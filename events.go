package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

//public class TabOpened
//{
//public Guid Id;
//public int TableNumber;
//public string Waiter;
//}

type TabOpened struct {
	Guid        uuid.UUID
	TableNumber int
	WaitStaff   string
}

func (t TabOpened) ToJson() []byte {
	j, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("json.Marshal: %s", err)
	}
	return j
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
		Guid:        guid,
		TableNumber: tableNumber,
		WaitStaff:   waitStaff,
	}
}

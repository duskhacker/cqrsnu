package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

//public class OpenTab
//{
//public Guid Id;
//public int TableNumber;
//public string Waiter;
//}

type OpenTab struct {
	Guid        uuid.UUID
	TableNumber int
	WaitStaff   string
}

func NewOpenTab(tableNumber int, waiter string) OpenTab {
	return OpenTab{
		Guid:        uuid.NewRandom(),
		TableNumber: tableNumber,
		WaitStaff:   waiter,
	}
}

func (o OpenTab) ToJson() []byte {
	j, err := json.Marshal(o)
	if err != nil {
		log.Fatalf("json.Marshal: %s", err)
	}
	return j
}

func (o OpenTab) FromJson(data []byte) OpenTab {
	var err error
	err = json.Unmarshal(data, &o)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return o
}

//func (o OpenTab) FromJson(json []byte) OpenTab {
//	var err error
//	err = json.Unmarshal(json, &o)
//
//}

package main

import "code.google.com/p/go-uuid/uuid"

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

func NewTabOpened(guid uuid.UUID, tableNumber int, waitStaff string) TabOpened {
	return TabOpened{
		Guid:        guid,
		TableNumber: tableNumber,
		WaitStaff:   waitStaff,
	}
}

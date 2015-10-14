package main

import (
	"encoding/json"
	"log"

	"code.google.com/p/go-uuid/uuid"
)

type CommandException struct {
	Guid    uuid.UUID
	Type    string
	Message string
}

func (c CommandException) FromJson(data []byte) CommandException {
	var err error
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return c
}

func NewCommandException(guid uuid.UUID, t string, msg string) CommandException {
	return CommandException{Guid: guid, Type: t, Message: msg}
}

func (c CommandException) Error() string {
	return c.Type + ":" + c.Message + "(" + c.Guid.String() + ")"
}

var tabNotOpenException = NewCommandException(nil, "TabNotOpen", "Cannot Place order without open Tab")
var drinksNotOutstanding = NewCommandException(nil, "DrinksNotOutstanding", "Cannot serve unordered drinks")

package main

import (
	"encoding/json"
	"log"
)

type CommandException struct {
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

func NewCommandException(t string, msg string) CommandException {
	return CommandException{Type: t, Message: msg}
}

func (c CommandException) Error() string {
	return c.Type + ":" + c.Message
}

var TabNotOpenException = NewCommandException("TabNotOpen", "Cannot Place order without open Tab")
var DrinksNotOutstanding = NewCommandException("DrinksNotOutstanding", "Cannot serve unordered drinks")

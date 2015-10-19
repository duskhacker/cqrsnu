package chef_todo_list

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/duskhacker/cqrsnu/cafe"
)

var chefTodoList = []*todoListGroup{}

type todoListItem struct {
	MenuNumber  int
	Description string
}

type todoListGroup struct {
	TabID uuid.UUID
	Items []todoListItem
}

func getTodoListGroup(tabID uuid.UUID) *todoListGroup {
	for _, list := range chefTodoList {
		if list.TabID.String() == tabID.String() {
			return list
		}
	}
	return nil
}

func newTodoListGroup(tabID uuid.UUID, items []cafe.OrderedItem) *todoListGroup {
	group := todoListGroup{TabID: tabID}
	for _, item := range items {
		group.Items = append(group.Items, todoListItem{MenuNumber: item.MenuNumber, Description: item.Description})
	}
	return &group
}

package chef_todo_list

import "sync"

var (
	mutex sync.RWMutex
)

package mtdmodels

import (
	"fmt"
	"os"
)

type Priority = int8

type ToDoItem struct {
	Id       int      `json:"id"`
	Summary  string   `json:"summary"`
	Done     bool     `json:"done,omitempty"`
	Priority Priority `json:"priority,omitempty"`
}

type ToDoGlobal = []ToDoItem

const PRIORITY_NOT_SET int8 = 0
const PRIORITY_HIGHT int8 = 1
const PRIORITY_MEDIUM int8 = 2
const PRIORITY_LOW int8 = 3

const MODE_EDIT = os.O_CREATE | os.O_RDWR
const MODE_READ = os.O_CREATE | os.O_RDONLY

type TodoListStorage interface {
	ReadTodoList() (error, *ToDoGlobal)
	SaveToDoList(lst *ToDoGlobal) error
	UseSource(sourceStr string)
}

type TodoListManager interface {
	AddItem(item string, priority Priority) error
	List() (error, *ToDoGlobal)
	Done(id int) error
	UseList(listName string)
}

type MtdError struct {
	Where         string
	Why           string
	OriginalError *error
}

func (e *MtdError) Error() string {
	return fmt.Sprintf("Error: %v, because of %s", e.Where, e.Why)
}

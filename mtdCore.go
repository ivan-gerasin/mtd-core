package mtdCore

import (
	"fmt"
	"log"
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

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const PRIORITY_NOT_SET int8 = 0
const PRIORITY_HIGHT int8 = 1
const PRIORITY_MEDIUM int8 = 2
const PRIORITY_LOW int8 = 3

const MODE_EDIT = os.O_CREATE | os.O_RDWR
const MODE_READ = os.O_CREATE | os.O_RDONLY

type TodoListStorage interface {
	ReadTodoList() (error, *ToDoGlobal)
	SaveToDoList(lst *ToDoGlobal) error
}

var Storage TodoListStorage = &FileStorage{}

type MtdError struct {
	Where         string
	Why           string
	OriginalError *error
}

func (e *MtdError) Error() string {
	return fmt.Sprintf("Error: %v, because of %s", e.Where, e.Why)
}

func AddItem(item string, priority Priority) error {
	err, ptrResults := Storage.ReadTodoList()
	if err != nil {
		return &MtdError{Why: err.Error(), Where: "AddItem(): While trying to read todo list", OriginalError: &err}
	}

	highestNumber := 0
	for key := range *ptrResults {
		if (*ptrResults)[key].Id > highestNumber {
			highestNumber = (*ptrResults)[key].Id
		}
	}
	*ptrResults = append(*ptrResults, ToDoItem{highestNumber + 1, item, false, priority})
	err = Storage.SaveToDoList(ptrResults)
	if err != nil {
		return &MtdError{Why: err.Error(), Where: "AddItem(): While trying to save todo list", OriginalError: &err}
	}
	return nil
}

func List() (error, *ToDoGlobal) {
	err, ptrResults := Storage.ReadTodoList()
	if err != nil {
		return &MtdError{Where: "List(): while reading todo list", Why: err.Error(), OriginalError: &err}, nil
	}
	return nil, ptrResults
}

func Done(id int) error {
	err, ptrResults := Storage.ReadTodoList()
	if id <= 0 {
		return &MtdError{Why: "number can not be equal to zero or lower", Where: "Done(): checking if id is valid", OriginalError: nil}
	}
	found := false
	for i := range *ptrResults {
		if (*ptrResults)[i].Id == id {
			(*ptrResults)[i].Done = true
			found = true
			break
		}
	}
	if found {
		err = Storage.SaveToDoList(ptrResults)
		if err != nil {
			return &MtdError{Why: err.Error(), Where: "Done(): While trying to save todo list after setting Done flag", OriginalError: &err}
		}
	} else {
		return &MtdError{Why: "no element with id#" + fmt.Sprint(id), Where: "Done(): while looking for given Id", OriginalError: nil}
	}
	return nil
}

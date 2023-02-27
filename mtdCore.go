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

func AddItem(item string, priority Priority) {
	err, ptrResults := Storage.ReadTodoList()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

func List() *ToDoGlobal {
	err, ptrResults := Storage.ReadTodoList()
	if err != nil {
		log.Fatal()
	}
	return ptrResults
}

func Done(id int) {
	err, ptrResults := Storage.ReadTodoList()
	if id <= 0 {
		fmt.Println("Error: number can not be equal to zero or lower") // TODO throw error
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
			log.Fatal("Failed to save todo list")
		}
	} else {
		fmt.Println("Error: no such element") // TODO throw error
	}
}

package mtdCore

import (
	"encoding/json"
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

func readTodoList(mode int) (*os.File, *ToDoGlobal, func()) {
	file, err := os.OpenFile("todolist.json", mode, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	size := fileStat.Size()
	buffer := make([]byte, size)
	readSize, err := file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	if int64(readSize) != size {
		log.Fatal("Read size and actual size are different")
	}
	if readSize == 0 {
		buffer = []byte(`[]`)
	}

	results := make(ToDoGlobal, 10) // TODO: figure out what is best way identify size
	err = json.Unmarshal(buffer, &results)
	if err != nil {
		log.Fatal(err)
	}

	closeFile := func() {
		file.Close()
	}

	return file, &results, closeFile
}

func saveToDoList(file *os.File, todoList *ToDoGlobal) {
	bytesToWrite, err := json.Marshal(*todoList)
	errCheck(err)

	_, err = file.Seek(0, 0)
	errCheck(err)
	_, err = file.Write(bytesToWrite)
	errCheck(err)
	file.Close()
}

func AddItem(item string, priority Priority) {
	file, ptrResults, _ := readTodoList(MODE_EDIT)
	highestNumber := 1
	for key := range *ptrResults {
		if (*ptrResults)[key].Id > highestNumber {
			highestNumber = (*ptrResults)[key].Id
		}
	}
	*ptrResults = append(*ptrResults, ToDoItem{highestNumber + 1, item, false, priority})
	saveToDoList(file, ptrResults)
}

func List() *ToDoGlobal {
	_, ptrResults, closeFile := readTodoList(MODE_READ)
	defer closeFile()

	return ptrResults
}

func Done(id int) {
	file, ptrResults, _ := readTodoList(MODE_EDIT)
	if id <= 0 {
		fmt.Println("Error: number can not be equal to zero or lower") // TODO throw error
	}
	if len(*ptrResults) >= id {
		for i := range *ptrResults {
			if (*ptrResults)[i].Id == id {
				(*ptrResults)[i].Done = true
				break
			}
		}

		saveToDoList(file, ptrResults)
	} else {
		fmt.Println("Error: no such element") // TODO throw error
	}
}

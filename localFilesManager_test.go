package mtdCore

import (
	"github.com/ivan-gerasin/mtdcore/mtdmodels"
	"testing"
)

type MockedStorage struct {
}

const testItemId = 99

func (fs MockedStorage) ReadTodoList() (error, *mtdmodels.ToDoGlobal) {
	return nil, &mtdmodels.ToDoGlobal{
		mtdmodels.ToDoItem{
			Id:       testItemId,
			Summary:  "Test item",
			Done:     false,
			Priority: 9,
		},
	}
}

func (fs MockedStorage) SaveToDoList(lst *mtdmodels.ToDoGlobal) error {
	saveToDoListCalls = lst
	return nil
}

var storageMock mtdmodels.TodoListStorage = MockedStorage{}
var saveToDoListCalls *mtdmodels.ToDoGlobal
var manager = LocalFilesManager{&storageMock}

func TestList(t *testing.T) {
	t.Run("Should return list of items", func(t *testing.T) {
		err, result := manager.List() // Expected list of 1 element
		if err != nil {
			t.Error("Error is not expected in this case")
		}
		if len(*result) != 1 {
			t.Error("Incorrect list size")
		}
		if (*result)[0].Id != testItemId {
			t.Error("Expecting different id value")
		}
		if (*result)[0].Summary != "Test item" {
			t.Error("Expecting different Summary value")
		}
		if (*result)[0].Done != false {
			t.Error("Expecting Done to be false")
		}
		if (*result)[0].Priority != 9 {
			t.Error("Expecting different Summary value")
		}
	})
}

func TestAddItem(t *testing.T) {
	t.Run("Should call SaveToDoList with new item", func(t *testing.T) {
		err := manager.AddItem("New text", 1)
		if err != nil {
			t.Error("Error is not expected in this case")
		}
		if saveToDoListCalls == nil {
			t.Error("Should not be nil")
		}
		if len(*saveToDoListCalls) != 2 {
			t.Error("Should add new element to list")
		}
		if (*saveToDoListCalls)[1].Id != testItemId+1 {
			t.Error("Should assign Id that greater than existing by 1")
		}
		if (*saveToDoListCalls)[1].Summary != "New text" {
			t.Error("Should save summary text")
		}
		if (*saveToDoListCalls)[1].Priority != 1 {
			t.Error("Should save priority value")
		}
		if (*saveToDoListCalls)[1].Done {
			t.Error("Should not set Done flag")
		}
	})
}

func TestDone(t *testing.T) {
	t.Run("Should mark given test item as done", func(t *testing.T) {
		err := manager.Done(testItemId)
		if err != nil {
			t.Error("Error is not expected in this case")
		}
		if (*saveToDoListCalls)[0].Done != true {
			t.Error("Should set Done flag")
		}
	})
}

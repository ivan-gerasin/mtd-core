package mtdCore

import (
	"fmt"
	"github.com/ivan-gerasin/mtdcore/mtdmodels"
)

type LocalFilesManager struct {
	storage *mtdmodels.TodoListStorage
}

func (manager *LocalFilesManager) AddItem(item string, priority mtdmodels.Priority) error {
	err, ptrResults := (*manager.storage).ReadTodoList()
	if err != nil {
		return &mtdmodels.MtdError{Why: err.Error(), Where: "AddItem(): While trying to read todo list", OriginalError: &err}
	}

	highestNumber := 0
	for key := range *ptrResults {
		if (*ptrResults)[key].Id > highestNumber {
			highestNumber = (*ptrResults)[key].Id
		}
	}
	*ptrResults = append(*ptrResults, mtdmodels.ToDoItem{highestNumber + 1, item, false, priority})
	err = (*manager.storage).SaveToDoList(ptrResults)
	if err != nil {
		return &mtdmodels.MtdError{Why: err.Error(), Where: "AddItem(): While trying to save todo list", OriginalError: &err}
	}
	return nil
}

func (manager *LocalFilesManager) List() (error, *mtdmodels.ToDoGlobal) {
	err, ptrResults := (*manager.storage).ReadTodoList()
	if err != nil {
		return &mtdmodels.MtdError{Where: "List(): while reading todo list", Why: err.Error(), OriginalError: &err}, nil
	}
	return nil, ptrResults
}

func (manager *LocalFilesManager) Done(id int) error {
	err, ptrResults := (*manager.storage).ReadTodoList()
	if id <= 0 {
		return &mtdmodels.MtdError{Why: "number can not be equal to zero or lower", Where: "Done(): checking if id is valid", OriginalError: nil}
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
		err = (*manager.storage).SaveToDoList(ptrResults)
		if err != nil {
			return &mtdmodels.MtdError{Why: err.Error(), Where: "Done(): While trying to save todo list after setting Done flag", OriginalError: &err}
		}
	} else {
		return &mtdmodels.MtdError{Why: "no element with id#" + fmt.Sprint(id), Where: "Done(): while looking for given Id", OriginalError: nil}
	}
	return nil
}

func (manager *LocalFilesManager) UseStorage(storage mtdmodels.TodoListStorage) {
	if storage == nil {
		panic("Storage can not be null")
	}
	manager.storage = &storage
}

func (manager *LocalFilesManager) UseList(listName string) {
	if listName != "" {
		(*manager.storage).UseSource(listName)
	} else {
		(*manager.storage).UseSource("default")
	}
}

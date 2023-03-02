package fileStorage

import (
	"encoding/json"
	"fmt"
	"github.com/ivan-gerasin/mtdcore/mtdmodels"
	"os"
)

type File interface {
	Close() error
	Seek(offset int64, whence int) (ret int64, err error)
	Stat() (os.FileInfo, error)
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

type FileSystem interface {
	OpenFile(name string, flag int, perm int) (File, error)
}

type StandardFileSystem struct{}

func (fs StandardFileSystem) OpenFile(name string, flag int, perm int) (File, error) {
	return os.OpenFile(name, flag, os.FileMode(perm))
}

var standardFileSystem = StandardFileSystem{}

type FileStorageError struct {
	Details       string
	OriginalError error
}

func (err *FileStorageError) Error() string {
	return fmt.Sprintf("FileStorageError: %s", err.Details)
}

func readTodoList(mode int) (error, File, *mtdmodels.ToDoGlobal, func() error) {
	file, err := standardFileSystem.OpenFile("todolist.json", mode, 0644)
	if err != nil {
		return &FileStorageError{"readTodoList(): Error while trying to open file", err},
			nil,
			nil,
			nil
	}

	fileStat, err := file.Stat()
	if err != nil {
		return &FileStorageError{"readTodoList(): Error while trying to get file Stat information", err},
			nil,
			nil,
			nil
	}

	size := fileStat.Size()
	buffer := make([]byte, size)
	readSize, err := file.Read(buffer)
	if err != nil {
		return &FileStorageError{"readTodoList(): Error while trying to read file", err},
			nil,
			nil,
			nil
	}
	if int64(readSize) != size {
		return &FileStorageError{"readTodoList(): Error there is a difference b/w file size by file.Stat() and read size", err},
			nil,
			nil,
			nil
	}
	if readSize == 0 {
		buffer = []byte(`[]`)
	}

	results := make(mtdmodels.ToDoGlobal, 10) // TODO: figure out what is best way identify size
	err = json.Unmarshal(buffer, &results)
	if err != nil {
		return &FileStorageError{"readTodoList(): fail to Unmarshal json file", err},
			nil,
			nil,
			nil
	}

	closeFile := func() error {
		err = file.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return nil, file, &results, closeFile
}

func saveToDoList(file File, todoList *mtdmodels.ToDoGlobal) error {
	bytesToWrite, err := json.Marshal(*todoList)
	if err != nil {
		return &FileStorageError{Details: "saveToDoList(): Error while json marshalling", OriginalError: err}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return &FileStorageError{Details: "saveToDoList(): Error while moving position in file", OriginalError: err}
	}

	_, err = file.Write(bytesToWrite)
	if err != nil {
		return &FileStorageError{Details: "saveToDoList(): Error while writing to file", OriginalError: err}
	}

	err = file.Close()
	if err != nil {
		return &FileStorageError{Details: "saveToDoList(): Error while closing the file", OriginalError: err}
	}
	return nil
}

type FileStorage struct{}

func (fs FileStorage) ReadTodoList() (error, *mtdmodels.ToDoGlobal) {
	err, _, list, closeFile := readTodoList(mtdmodels.MODE_READ)
	if err != nil {
		return &mtdmodels.MtdError{Where: "ReadTodoList(): failed to read todo list with FileStorage", Why: err.Error(), OriginalError: &err}, nil
	}
	err = closeFile()
	if err != nil {
		return &mtdmodels.MtdError{Where: "ReadTodoList(): failed to close file with FileStorage", Why: err.Error(), OriginalError: &err}, nil
	}
	return nil, list
}

func (fs FileStorage) SaveToDoList(lst *mtdmodels.ToDoGlobal) error {
	err, file, _, closeFile := readTodoList(mtdmodels.MODE_EDIT)
	if err != nil {
		return &mtdmodels.MtdError{Where: "SaveToDoList(): failed to read todo list with FileStorage", Why: err.Error(), OriginalError: &err}
	}
	err = saveToDoList(file, lst)
	if err != nil {
		return &mtdmodels.MtdError{Where: "SaveToDoList(): failed to save todo list file with FileStorage", Why: err.Error(), OriginalError: &err}
	}
	err = closeFile()
	if err != nil {
		return &mtdmodels.MtdError{Where: "SaveToDoList(): failed to close file with FileStorage", Why: err.Error(), OriginalError: &err}
	}
	return nil
}

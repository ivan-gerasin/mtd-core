package mtdCore

import (
	"encoding/json"
	"log"
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

var fs = StandardFileSystem{}

func readTodoList(mode int) (File, *ToDoGlobal, func()) {
	//file, err := os.OpenFile("todolist.json", mode, 0644)
	file, err := fs.OpenFile("todolist.json", mode, 0644)
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

func saveToDoList(file File, todoList *ToDoGlobal) {
	bytesToWrite, err := json.Marshal(*todoList)
	errCheck(err)

	_, err = file.Seek(0, 0)
	errCheck(err)
	_, err = file.Write(bytesToWrite)
	errCheck(err)
	file.Close()
}

type FileStorage struct{}

/*

	ONE GIANT TODO: RETURN ERRORS!

*/

func (fs FileStorage) ReadTodoList() (error, *ToDoGlobal) {
	_, list, closeFile := readTodoList(MODE_READ)
	closeFile()
	return nil, list
}

func (fs FileStorage) SaveToDoList(lst *ToDoGlobal) error {
	file, _, closeFile := readTodoList(MODE_EDIT)
	saveToDoList(file, lst)
	closeFile()
	return nil
}

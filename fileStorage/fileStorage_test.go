package fileStorage

import (
	"os"
)

type MockedFileSystem struct {
}
type FileMock struct {
	close error
	seek  struct {
		ret int64
		err error
	}
	stat struct {
		finfo os.FileInfo
		error
	}
	read struct {
		n   int
		err error
	}
	write struct {
		n   int
		err error
	}
}

func (fileMock FileMock) Close() error {
	return fileMock.close
}

func (fileMock FileMock) Seek(_ int64, _ int) (ret int64, err error) {
	return fileMock.seek.ret, fileMock.seek.err
}

func (fileMock FileMock) Stat() (os.FileInfo, error) {
	return fileMock.stat.finfo, fileMock.stat.error
}

func (fileMock FileMock) Read(_ []byte) (n int, err error) {
	return fileMock.read.n, fileMock.read.err
}

func (fileMock FileMock) Write(_ []byte) (n int, err error) {
	return fileMock.write.n, fileMock.write.err
}

var fileMockPrototype = FileMock{}

func (fs MockedFileSystem) OpenFile(name string, flag int, perm int) (File, error) {
	return fileMockPrototype, nil
}

var mockedFs = &MockedFileSystem{}
var originalFs FileSystem = nil

func setup() {
	originalFs = standardFileSystem
	standardFileSystem = mockedFs
}

func tearDown() {
	standardFileSystem = originalFs
}

// TODO: someday write tests for fileStorage :)

//func TestFileStorage_ReadTodoList(t *testing.T) {
//	t.Run("Access to a file within filesystem with OpenFile", func(t *testing.T) {
//		setup()
//
//		tearDown()
//	})
//}

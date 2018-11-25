package tester

import (
	"fmt"

	"github.com/pw1/stor"
	"github.com/stretchr/testify/suite"
)

// Generic test suite for implementations of the Storage interface
type StorageTester struct {
	suite.Suite

	storage     stor.Storage
	storageType stor.Type

	SetupSuiteFunc    func(*StorageTester)
	TearDownSuiteFunc func(*StorageTester)

	SetupTestFunc    func(*StorageTester) stor.Storage
	TearDownTestFunc func(*StorageTester)
}

func New(storageType stor.Type) *StorageTester {
	st := &StorageTester{
		storageType: storageType,
	}
	return st
}

func (s *StorageTester) SetupSuite() {
	if s.SetupSuiteFunc != nil {
		s.SetupSuiteFunc(s)
	}
}

func (s *StorageTester) TearDownSuite() {
	if s.TearDownSuiteFunc != nil {
		s.TearDownSuiteFunc(s)
	}
}

func (s *StorageTester) SetupTest() {
	if s.SetupTestFunc == nil {
		s.Fail("You have not defined the SetupTest function. This function is required.")
		return
	}
	s.storage = s.SetupTestFunc(s)
}

func (s *StorageTester) TearDownTest() {
	if s.TearDownTestFunc != nil {
		s.TearDownTestFunc(s)
	}
	s.storage = nil
}

func (s *StorageTester) insertStandardFiles() {
	files := map[string]string{
		"file1":           "test123",
		"dir1/file2":      "test456",
		"dir1/file3":      "test789",
		"dir1/dir4/file5": "test788909",
		"dir2/dir3/file4": "test0123",
	}

	for filepath, content := range files {
		err := s.storage.Save(filepath, []byte(content))
		if err != nil {
			msg := fmt.Sprintf("Failed to prepare test:\n  -> Failed to insert file: %s", filepath)
			msg += fmt.Sprintf("\n  -> Error: %s", err)
			s.Require().Fail(msg)
		}
	}
}

func (s *StorageTester) TestList() {
	s.insertStandardFiles()

	files, dirs, err := s.storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1", "dir2"}, dirs)
}

func (s *StorageTester) TestListEscapes() {
	files, dirs, err := s.storage.List("..")
	s.Empty(files)
	s.Empty(dirs)
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

func (s *StorageTester) TestListDot() {
	s.insertStandardFiles()

	files, dirs, err := s.storage.List(".")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1", "dir2"}, dirs)
}

func (s *StorageTester) TestListDir1() {
	s.insertStandardFiles()

	files, dirs, err := s.storage.List("dir1")
	s.Nil(err)
	s.ElementsMatch([]string{"dir1/file2", "dir1/file3"}, files)
	s.ElementsMatch([]string{"dir1/dir4"}, dirs)
}

func (s *StorageTester) TestExist() {
	s.insertStandardFiles()

	exist, err := s.storage.Exist("dir1/file2")
	s.True(exist)
	s.Nil(err)
}

func (s *StorageTester) TestExistNonExisting() {
	s.insertStandardFiles()

	exist, err := s.storage.Exist("dir45/file89")
	s.False(exist)
	s.Nil(err)
}

func (s *StorageTester) TestExistEscapes() {
	exist, err := s.storage.Exist("../outside.txt")
	s.False(exist)
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

func (s *StorageTester) TestLoad() {
	s.insertStandardFiles()

	data, err := s.storage.Load("file1", 1e6)
	s.Nil(err)
	s.Equal([]byte("test123"), data)
}

func (s *StorageTester) TestLoadEscapes() {
	s.insertStandardFiles()

	data, err := s.storage.Load("../file1", 1e6)
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
	s.Equal([]byte{}, data)
}

func (s *StorageTester) TestLoadInDir() {
	s.insertStandardFiles()

	data, err := s.storage.Load("dir1/file2", 1e6)
	s.Nil(err)
	s.Equal([]byte("test456"), data)
}

func (s *StorageTester) TestLoadWithMaxSize() {
	s.insertStandardFiles()

	data, err := s.storage.Load("file1", 6)
	s.NotNil(err)
	s.Equal([]byte{}, data)
}

func (s *StorageTester) TestLoadNonExisting() {
	s.insertStandardFiles()

	data, err := s.storage.Load("dir1/file1", 1e6)
	s.NotNil(err)
	s.Equal([]byte{}, data)
}

func (s *StorageTester) TestSave() {
	s.insertStandardFiles()

	testFile := "dir1/new-file.txt"
	testData := []byte("my-data")

	err := s.storage.Save(testFile, testData)
	s.Nil(err)

	savedData, err := s.storage.Load(testFile, 1e6)
	s.Nil(err)
	s.Equal(testData, savedData)
}

func (s *StorageTester) TestSaveOverwrite() {
	s.insertStandardFiles()

	testFile := "file1"
	testData := []byte("my-data")

	err := s.storage.Save(testFile, testData)
	s.Nil(err)

	savedData, err := s.storage.Load(testFile, 1e6)
	s.Nil(err)
	s.Equal(testData, savedData)
}

func (s *StorageTester) TestSaveEscapes() {
	s.insertStandardFiles()

	err := s.storage.Save("../file1", []byte("qwerty"))
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

func (s *StorageTester) TestDelete() {
	s.insertStandardFiles()

	err := s.storage.Delete("dir1/file2")
	s.Nil(err)

	_, err = s.storage.Load("dir1/file2", 1e6)
	s.NotNil(err)

	files, dirs, err := s.storage.List("dir1")
	s.Nil(err)
	s.ElementsMatch([]string{"dir1/file3"}, files)
	s.ElementsMatch([]string{"dir1/dir4"}, dirs)
}

func (s *StorageTester) TestDeleteDir() {
	s.insertStandardFiles()

	err := s.storage.Delete("dir2/dir3/file4")
	s.Nil(err)

	_, err = s.storage.Load("dir2/dir3/file4", 1e6)
	s.NotNil(err)

	files, dirs, err := s.storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1"}, dirs)
}

func (s *StorageTester) TestDeleteAll() {
	s.insertStandardFiles()

	err := s.storage.Delete("file1")
	s.Nil(err)
	err = s.storage.Delete("dir1/file2")
	s.Nil(err)
	err = s.storage.Delete("dir1/file3")
	s.Nil(err)
	err = s.storage.Delete("dir1/dir4/file5")
	s.Nil(err)
	err = s.storage.Delete("dir2/dir3/file4")
	s.Nil(err)

	files, dirs, err := s.storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{}, files)
	s.ElementsMatch([]string{}, dirs)
}

func (s *StorageTester) TestDeleteEscapes() {
	s.insertStandardFiles()

	err := s.storage.Delete("../file1")
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

func (s *StorageTester) TestType() {
	s.Equal(s.storageType, s.storage.Type())
}

// Package tester contains a test suite to test whether objects correctly implement the stor.Storage
// interface.
//
// Example usage:
//
//  import (
// 	    "testing"
// 	    "github.com/stretchr/testify/suite"
// 	    "github.com/pw1/stor"
// 	    "github.com/pw1/stor/tester"
//  )
//
//  func TestMyStorageTester(t *testing.T) {
//      myConfFactory := func() *stor.Conf {
//          return &stor.Conf{
//              Type: MyStorageType,
//          }
//      }
//
// 	    testSuite := &tester.StorageTester{
// 		    ConfFactory: myConfFactory,
// 	    }
//
// 	    suite.Run(t, testSuite)
//  }
//
package tester

import (
	"fmt"

	"github.com/pw1/stor"
	"github.com/stretchr/testify/suite"
)

// StorageTester is generic test suite for implementations of the Storage interface
type StorageTester struct {
	suite.Suite

	// Storage is the storage object that is tested. All tests expect the storage object under test
	// in this variable.
	// If the ConfFactory is defined, then this is set before each test by calling this factory to
	// create a configuration, and then calling stor.New with that configuration. Alternatively, you
	// can define a SetupTestFunc that creates a new storage object and saves it in the
	// StorageTester.Storage variable.
	Storage stor.Storage

	// ConfFactory is the factory function for creating
	ConfFactory func() *stor.Conf

	// SetupSuiteFunc is the function that is called once before the first test is run.
	SetupSuiteFunc func(*StorageTester)

	// SetupSuiteFunc is the function that is called once after all tests are executed.
	TearDownSuiteFunc func(*StorageTester)

	// SetupTestFunc is called before before each test.
	SetupTestFunc func(*StorageTester)

	// TearDownTestFunc is called after each test.
	TearDownTestFunc func(*StorageTester)
}

// SetupSuite is executed before the first test is executed. It will call SetupSuiteFunc if that is
// defined.
func (s *StorageTester) SetupSuite() {
	if s.SetupSuiteFunc != nil {
		s.SetupSuiteFunc(s)
	}
}

// TearDownSuite is called after the last test is executed. It will call TearDownSuiteFunc if that
// is defined.
func (s *StorageTester) TearDownSuite() {
	if s.TearDownSuiteFunc != nil {
		s.TearDownSuiteFunc(s)
	}
}

// SetupTest is called before each test is executed. If s.ConfFactory is defined, then it will call
// that function to create a new configuration, and then call stor.New() with that new
// configuration. The resulting stor.Storage is saved to s.Storage. It will SetupTestFunc if that is
// defined.
func (s *StorageTester) SetupTest() {
	if s.ConfFactory != nil {
		st, err := stor.New(s.ConfFactory())
		if err != nil {
			s.FailNow("failed to create new Storage object", err)
		}
		s.Storage = st
	}

	if s.SetupTestFunc != nil {
		s.SetupTestFunc(s)
	}
}

// TearDownTest is called before each test is executed. It will execute TearDownTestFunc is that is
// defined. It will also set s.Storage to nil (the Storage must be recreated before each test).
func (s *StorageTester) TearDownTest() {
	if s.TearDownTestFunc != nil {
		s.TearDownTestFunc(s)
	}
	s.Storage = nil
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
		err := s.Storage.Save(filepath, []byte(content))
		if err != nil {
			msg := fmt.Sprintf("Failed to prepare test:\n  -> Failed to insert file: %s", filepath)
			msg += fmt.Sprintf("\n  -> Error: %s", err)
			s.Require().Fail(msg)
		}
	}
}

// TestMeta verifies that Meta() returns meta information about a file.
func (s *StorageTester) TestMeta() {
	s.insertStandardFiles()

	meta, err := s.Storage.Meta("dir1/file3")
	s.Nil(err)
	s.Equal(&stor.Meta{
		Size: 7,
	}, meta)
}

// TestMetaEscapes verifies that Meta() returns an error if the supplied path is invalid.
func (s *StorageTester) TestMetaEscapes() {
	s.insertStandardFiles()

	meta, err := s.Storage.Meta("../file1")
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
	s.Nil(meta)
}

// TestMetaNonExisting verifies that Meta() returns an error if the supplied path doesn't exist.
func (s *StorageTester) TestMetaNonExisting() {
	s.insertStandardFiles()

	meta, err := s.Storage.Meta("dir1/file1")
	s.NotNil(err)
	s.True(stor.IsPathDoesntExistError(err))
	s.Nil(meta)
}

// TestList verifies that List() returns a list of files and subdirectories in the root of the
// storage.
func (s *StorageTester) TestList() {
	s.insertStandardFiles()

	files, dirs, err := s.Storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1", "dir2"}, dirs)
}

// TestListEscapes verifies that List() returns an error if the supplied path is invalid.
func (s *StorageTester) TestListEscapes() {
	files, dirs, err := s.Storage.List("..")
	s.Empty(files)
	s.Empty(dirs)
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

// TestListDot verifies that List(".") lists the files and subdirectories in the root of the
// storage.
func (s *StorageTester) TestListDot() {
	s.insertStandardFiles()

	files, dirs, err := s.Storage.List(".")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1", "dir2"}, dirs)
}

// TestListDir1 verifies that List() returns files and subdirectories in a directory.
func (s *StorageTester) TestListDir1() {
	s.insertStandardFiles()

	files, dirs, err := s.Storage.List("dir1")
	s.Nil(err)
	s.ElementsMatch([]string{"dir1/file2", "dir1/file3"}, files)
	s.ElementsMatch([]string{"dir1/dir4"}, dirs)
}

// TestLoad verifies that Load() returns the content of a file.
func (s *StorageTester) TestLoad() {
	s.insertStandardFiles()

	data, err := s.Storage.Load("file1", 1e6)
	s.Nil(err)
	s.Equal([]byte("test123"), data)
}

// TestLoadEscapes verifies that Load() returns an error if the supplied path is invalid.
func (s *StorageTester) TestLoadEscapes() {
	s.insertStandardFiles()

	data, err := s.Storage.Load("../file1", 1e6)
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
	s.Equal([]byte{}, data)
}

// TestLoadInDir verifies that Load() returns the content of a file in a directory.
func (s *StorageTester) TestLoadInDir() {
	s.insertStandardFiles()

	data, err := s.Storage.Load("dir1/file2", 1e6)
	s.Nil(err)
	s.Equal([]byte("test456"), data)
}

// TestLoadWithMaxSize verifies that Load() returns an error if the specified file is larger than
// the specified maximum size.
func (s *StorageTester) TestLoadWithMaxSize() {
	s.insertStandardFiles()

	data, err := s.Storage.Load("file1", 6)
	s.NotNil(err)
	s.True(stor.IsTooLargeError(err))
	s.Equal([]byte{}, data)
}

// TestLoadNonExisting verifies that Load() returns an error if the supplied path doesn't exist.
func (s *StorageTester) TestLoadNonExisting() {
	s.insertStandardFiles()

	data, err := s.Storage.Load("dir1/file1", 1e6)
	s.NotNil(err)
	s.True(stor.IsPathDoesntExistError(err))
	s.Equal([]byte{}, data)
}

// TestSave verifies that Save() saves data to a file.
func (s *StorageTester) TestSave() {
	s.insertStandardFiles()

	testFile := "dir1/new-file.txt"
	testData := []byte("my-data")

	err := s.Storage.Save(testFile, testData)
	s.Nil(err)

	savedData, err := s.Storage.Load(testFile, 1e6)
	s.Nil(err)
	s.Equal(testData, savedData)
}

// TestSaveOverwrite verifies that Save() overwrites an existing file without any error.
func (s *StorageTester) TestSaveOverwrite() {
	s.insertStandardFiles()

	testFile := "file1"
	testData := []byte("my-data")

	err := s.Storage.Save(testFile, testData)
	s.Nil(err)

	savedData, err := s.Storage.Load(testFile, 1e6)
	s.Nil(err)
	s.Equal(testData, savedData)
}

// TestSaveEscapes verifies that Save() returns an error if the supplied path is invalid.
func (s *StorageTester) TestSaveEscapes() {
	s.insertStandardFiles()

	err := s.Storage.Save("../file1", []byte("qwerty"))
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

// TestDelete verifies that Delete() removes a file from storage.
func (s *StorageTester) TestDelete() {
	s.insertStandardFiles()

	err := s.Storage.Delete("dir1/file2")
	s.Nil(err)

	_, err = s.Storage.Load("dir1/file2", 1e6)
	s.NotNil(err)

	files, dirs, err := s.Storage.List("dir1")
	s.Nil(err)
	s.ElementsMatch([]string{"dir1/file3"}, files)
	s.ElementsMatch([]string{"dir1/dir4"}, dirs)
}

// TestDeleteDir verifies that if the last file inside a subdirectory is removed, that the parent
// subdirectory (which is now empty) is also removed.
func (s *StorageTester) TestDeleteDir() {
	s.insertStandardFiles()

	err := s.Storage.Delete("dir2/dir3/file4")
	s.Nil(err)

	_, err = s.Storage.Load("dir2/dir3/file4", 1e6)
	s.NotNil(err)

	files, dirs, err := s.Storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{"file1"}, files)
	s.ElementsMatch([]string{"dir1"}, dirs)
}

// TestDeleteNonExisting verifies that Delete() returns an error if the supplied path doesn't exist.
func (s *StorageTester) TestDeleteNonExisting() {
	err := s.Storage.Delete("dir1/file1")
	s.NotNil(err)
	s.True(stor.IsPathDoesntExistError(err))
}

// TestDeleteAll verifies if all files are deleted one by one that the storage is empty afterwards.
func (s *StorageTester) TestDeleteAll() {
	s.insertStandardFiles()

	err := s.Storage.Delete("file1")
	s.Nil(err)
	err = s.Storage.Delete("dir1/file2")
	s.Nil(err)
	err = s.Storage.Delete("dir1/file3")
	s.Nil(err)
	err = s.Storage.Delete("dir1/dir4/file5")
	s.Nil(err)
	err = s.Storage.Delete("dir2/dir3/file4")
	s.Nil(err)

	files, dirs, err := s.Storage.List("")
	s.Nil(err)
	s.ElementsMatch([]string{}, files)
	s.ElementsMatch([]string{}, dirs)
}

// TestDeleteEscapes verifies that Delete() returns an error if the supplied path is invalid.
func (s *StorageTester) TestDeleteEscapes() {
	s.insertStandardFiles()

	err := s.Storage.Delete("../file1")
	s.NotNil(err)
	s.True(stor.IsInvalidPathError(err))
}

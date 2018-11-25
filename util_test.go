package stor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestStorageUtil(t *testing.T) {
	suite.Run(t, new(StorageUtilSuite))
}

type StorageUtilSuite struct {
	suite.Suite
}

// Test valid clean paths.
func (s *StorageUtilSuite) TestCleanPath() {
	table := [][]string{
		[]string{"file1", "file1"},
		[]string{"fiLe1/", "fiLe1"},
		[]string{"dir1/file.1", "dir1/file.1"},
		[]string{"dir1/file-1//", "dir1/file-1"},
		[]string{"dir1//file_1", "dir1/file_1"},
		[]string{"", ""},
		[]string{".", ""},
		[]string{"./", ""},
	}

	for _, row := range table {
		inputPath := row[0]
		expectedPath := row[1]

		cleanPath, err := CleanPath(inputPath)
		msg := fmt.Sprintf("Input: %s, Expected output: %s, Actual output: %s",
			inputPath, expectedPath, cleanPath)
		s.Equal(expectedPath, cleanPath, msg)
		s.Nil(err)
	}
}

// Test invalid paths. All these path must return an error
func (s *StorageUtilSuite) TestCleanPathInvalid() {
	table := []string{
		"/absolute/file1",
		"../file1",
		"dir1/../file1",
		"file*1",
		"dir1/filé1",
		"dir1\\file1",
		"c:\\dir1\\file1",
		"D:/dir1/file1",
	}

	for _, inputPath := range table {
		cleanPath, err := CleanPath(inputPath)
		msg := fmt.Sprintf("Input: %s, Expected output is Empty, Actual output: %s",
			inputPath, cleanPath)
		s.Empty(cleanPath, msg)
		s.NotNil(err)
		s.True(IsInvalidPathError(err), fmt.Sprintf("Input: %s, Actual error: %v", inputPath, err))
	}
}

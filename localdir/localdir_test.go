package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pw1/stor"
	"github.com/pw1/stor/tester"
)

// Call the generic storage tests
func TestLocalDirWithStorageTester(t *testing.T) {
	var tempDir string
	tempDir, err := ioutil.TempDir("", "TestLocalDirGenericStorage")
	if err != nil {
		t.FailNow()
	}
	t.Logf("Temp dir for testing: %s", tempDir)

	myConfFactory := func() *stor.Conf {
		return &stor.Conf{
			Type: LocalDirStorageType,
			Path: tempDir,
		}
	}

	testSuite := &tester.StorageTester{
		ConfFactory:       myConfFactory,
		SetupTestFunc:     func(s *tester.StorageTester) { cleanDir(t, tempDir) },
		TearDownSuiteFunc: func(s *tester.StorageTester) { os.RemoveAll(tempDir) },
	}
	suite.Run(t, testSuite)
}

// cleanDir removes all files and subdirectories. But it does not remove the directory itself.
func cleanDir(t *testing.T, dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		t.Fatalf("Failed to list dir: %s", err)
	}

	for _, file := range files {
		err = os.RemoveAll(filepath.Join(dirPath, file.Name()))
		if err != nil {
			t.Fatalf("Failed to cleanup dir: %s", err)
		}
	}
}

// Create a new directory to be used by a single test
func makeTestDir(tempDir string) (string, error) {
	tempDir, err := ioutil.TempDir(tempDir, "")
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

// LocalDirSuite contains tests that are not in the generic tester.StorageTester suite. The tests in
// this suite are specific for the LocalDir type.
type LocalDirSuite struct {
	suite.Suite
	tempDir string
}

func TestLocalDirSuite(t *testing.T) {
	suite.Run(t, new(LocalDirSuite))
}

func (s *LocalDirSuite) SetupSuite() {
	newDir, err := ioutil.TempDir("", "TestSuiteLocalDir")
	s.Nil(err)
	s.tempDir = newDir
}

func (s *LocalDirSuite) TearDownSuite() {
	os.RemoveAll(s.tempDir)
}

func (s *LocalDirSuite) TestNewLocalDirAbs() {
	testDir, err := makeTestDir(s.tempDir)
	s.Nil(err)

	stConf := &stor.Conf{
		Type: LocalDirStorageType,
		Path: testDir,
	}

	localDir, err := New(stConf)
	s.Nil(err)
	s.NotNil(localDir)
	s.Equal(testDir, localDir.BaseDir)
}

func (s *LocalDirSuite) TestNewLocalDirRel() {
	testDir, err := makeTestDir(s.tempDir)
	s.Nil(err)

	myBaseDir := filepath.Join(testDir, "base")
	os.Mkdir(myBaseDir, 0700)

	os.Chdir(testDir)

	stConf := &stor.Conf{
		Type: LocalDirStorageType,
		Path: "base",
	}

	localDir, err := New(stConf)
	s.Nil(err)
	s.NotNil(localDir)
	s.Equal(myBaseDir, localDir.BaseDir)
}

func (s *LocalDirSuite) TestNewLocalDirNonExisting() {
	stConf := &stor.Conf{
		Type: LocalDirStorageType,
		Path: "_this_directory_doesnt_exist__",
	}

	localDir, err := New(stConf)
	s.NotNil(err)
	s.Nil(localDir)
}

// TestNewLocalDirFileBase verifies that that New() doesn't accept a file as BaseDir
func (s *LocalDirSuite) TestNewLocalDirFileBase() {
	testDir, err := makeTestDir(s.tempDir)
	s.Nil(err)

	myBaseFile := filepath.Join(testDir, "base")
	ioutil.WriteFile(myBaseFile, []byte("test123"), 0600)

	stConf := &stor.Conf{
		Type: LocalDirStorageType,
		Path: myBaseFile,
	}

	localDir, err := New(stConf)
	s.NotNil(err)
	s.Nil(localDir)
}

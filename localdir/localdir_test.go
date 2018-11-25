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

// Create a new directory to be used by a single test
func makeTestDir(tempDir string) (string, error) {
	tempDir, err := ioutil.TempDir(tempDir, "")
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

// Create a base directory (for the LocalDir storage) and the associated LocalDir object
func makeLocalDir(tempDir string) (*LocalDir, error) {
	testDir, err := makeTestDir(tempDir)
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(testDir, "base")
	os.Mkdir(baseDir, 0700)

	stConf := stor.NewConf()
	stConf.StorageType = LocalDirStorageType
	stConf.Path = baseDir

	localDir, err := New(stConf)
	if err != nil {
		return nil, err
	}

	return localDir, nil
}

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

	stConf := stor.NewConf()
	stConf.StorageType = LocalDirStorageType
	stConf.Path = testDir

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

	stConf := stor.NewConf()
	stConf.StorageType = LocalDirStorageType
	stConf.Path = "base"

	localDir, err := New(stConf)
	s.Nil(err)
	s.NotNil(localDir)
	s.Equal(myBaseDir, localDir.BaseDir)
}

func (s *LocalDirSuite) TestNewLocalDirNonExisting() {
	stConf := stor.NewConf()
	stConf.StorageType = LocalDirStorageType
	stConf.Path = "_this_directory_doesnt_exist__"

	localDir, err := New(stConf)
	s.NotNil(err)
	s.Nil(localDir)
}

// Test that New() doesn't accept a file as BaseDir
func (s *LocalDirSuite) TestNewLocalDirFileBase() {
	testDir, err := makeTestDir(s.tempDir)
	s.Nil(err)

	myBaseFile := filepath.Join(testDir, "base")
	ioutil.WriteFile(myBaseFile, []byte("test123"), 0600)

	stConf := stor.NewConf()
	stConf.StorageType = LocalDirStorageType
	stConf.Path = myBaseFile

	localDir, err := New(stConf)
	s.NotNil(err)
	s.Nil(localDir)
}

// Call the generic storage tests
func TestLocalDirGenericStorage(t *testing.T) {
	testSuite := tester.New(LocalDirStorageType)

	var tempDir string

	testSuite.SetupSuiteFunc = func(s *tester.StorageTester) {
		var err error
		tempDir, err = ioutil.TempDir("", "TestLocalDirGenericStorage")
		s.Nil(err)
	}

	testSuite.TearDownSuiteFunc = func(s *tester.StorageTester) {
		os.RemoveAll(tempDir)
	}

	testSuite.SetupTestFunc = func(s *tester.StorageTester) stor.Storage {
		ldir, err := makeLocalDir(tempDir)
		s.Nil(err)
		return ldir
	}

	suite.Run(t, testSuite)
}

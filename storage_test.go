package stor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

//
// Tests that starte the actual test suites in this file
//

func TestStorageTypeSuite(t *testing.T) {
	suite.Run(t, new(StorageTypeSuite))
}

func TestStorageErrorsSuite(t *testing.T) {
	suite.Run(t, new(StorageErrorsSuite))
}

//
// Test suite for the error types in Storage
//
type StorageErrorsSuite struct {
	suite.Suite
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorTrue() {
	err := NewInvalidPathError("123456")
	s.True(IsInvalidPathError(err))
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorFalse() {
	s.False(IsInvalidPathError(fmt.Errorf("Some other error")))
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorNil() {
	s.False(IsInvalidPathError(nil))
}

//
// Test Suite for the StorageType
//
type StorageTypeSuite struct {
	suite.Suite
	storageType      Type
	storageTypeText  string
	storageType2     Type
	storageTypeText2 string
}

func (s *StorageTypeSuite) SetupSuite() {
	s.storageType = 123
	s.storageTypeText = "MyTestingType"
	s.storageType2 = 456
	s.storageTypeText2 = "TheOtherTestingType"
	RegisterStorageType(s.storageType, s.storageTypeText)
	RegisterStorageType(s.storageType2, s.storageTypeText2)
}

func (s *StorageTypeSuite) TestConf() {
	s.NotNil(NewConf())
}

func (s *StorageTypeSuite) TestStorageTypeFmtString() {
	s.Equal(s.storageTypeText, fmt.Sprintf("%s", s.storageType))
}

func (s *StorageTypeSuite) TestStorageTypeFmtV() {
	s.Equal(s.storageTypeText, fmt.Sprintf("%v", s.storageType))
}

func (s *StorageTypeSuite) TestStorageTypeFmtInt() {
	s.Equal("456", fmt.Sprintf("%d", s.storageType2))
}

func (s *StorageTypeSuite) TestStorageTypeToInt() {
	s.Equal(456, int(s.storageType2))
}

func (s *StorageTypeSuite) TestStorageTypeUnmarshal() {
	var typ Type
	err := typ.UnmarshalText([]byte(s.storageTypeText))
	s.Nil(err)
	s.Equal(s.storageType, typ)
}

func (s *StorageTypeSuite) TestRegisterStorageType() {
	s.Nil(RegisterStorageType(123456789, "MyNewType"))
}

func (s *StorageTypeSuite) TestRegisterStorageTypeDuplInt() {
	s.NotNil(RegisterStorageType(2, "MyNewType"))
}
func (s *StorageTypeSuite) TestRegisterStorageTypeDuplStr() {
	s.NotNil(RegisterStorageType(789, s.storageTypeText))
}

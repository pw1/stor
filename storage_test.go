package stor

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

//
// Tests that starte the actual test suites in this file
//

func TestTypeSuite(t *testing.T) {
	suite.Run(t, new(TypeSuite))
}

func TestStorageErrorsSuite(t *testing.T) {
	suite.Run(t, new(StorageErrorsSuite))
}

func TestNewSuite(t *testing.T) {
	suite.Run(t, new(NewSuite))
}

//
// Test suite for the error types in Storage
//
type StorageErrorsSuite struct {
	suite.Suite
}

func (s *StorageErrorsSuite) TestUnregisteredTypeError() {
	err := &UnregisteredTypeError{"bla"}
	s.Contains(err.Error(), "bla")
}

func (s *StorageErrorsSuite) TestInvalidPathError() {
	err := &InvalidPathError{Path: "qwerty", Msg: "zxcvbn"}
	s.Contains(err.Error(), "qwerty")
	s.Contains(err.Error(), "zxcvbn")
}

func (s *StorageErrorsSuite) TestInvalidPathErrorNoMsg() {
	err := &InvalidPathError{Path: "qwerty"}
	s.Contains(err.Error(), "qwerty")
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorTrue() {
	err := &InvalidPathError{Path: "123456"}
	s.True(IsInvalidPathError(err))
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorFalse() {
	err := fmt.Errorf("Some other error")
	s.False(IsInvalidPathError(err))
}

func (s *StorageErrorsSuite) TestIsInvalidPathErrorNil() {
	s.False(IsInvalidPathError(nil))
}

func (s *StorageErrorsSuite) TestIsUnregisteredTypeError() {
	s.True(IsUnregisteredTypeError(&UnregisteredTypeError{}))
	s.False(IsUnregisteredTypeError(&InvalidPathError{}))
	s.False(IsUnregisteredTypeError(&PathDoesntExistError{}))
	s.False(IsUnregisteredTypeError(&TooLargeError{}))
	s.False(IsUnregisteredTypeError(&UnspecifiedTypeError{}))
	s.False(IsUnregisteredTypeError(errors.New("test")))
}

func (s *StorageErrorsSuite) TestIsInvalidPathError() {
	s.False(IsInvalidPathError(&UnregisteredTypeError{}))
	s.True(IsInvalidPathError(&InvalidPathError{}))
	s.False(IsInvalidPathError(&PathDoesntExistError{}))
	s.False(IsInvalidPathError(&TooLargeError{}))
	s.False(IsInvalidPathError(&UnspecifiedTypeError{}))
	s.False(IsInvalidPathError(errors.New("test")))
}

func (s *StorageErrorsSuite) TestIsPathDoesntExistError() {
	s.False(IsPathDoesntExistError(&UnregisteredTypeError{}))
	s.False(IsPathDoesntExistError(&InvalidPathError{}))
	s.True(IsPathDoesntExistError(&PathDoesntExistError{}))
	s.False(IsPathDoesntExistError(&TooLargeError{}))
	s.False(IsPathDoesntExistError(&UnspecifiedTypeError{}))
	s.False(IsPathDoesntExistError(errors.New("test")))
}

func (s *StorageErrorsSuite) TestIsTooLargeError() {
	s.False(IsTooLargeError(&UnregisteredTypeError{}))
	s.False(IsTooLargeError(&InvalidPathError{}))
	s.False(IsTooLargeError(&PathDoesntExistError{}))
	s.True(IsTooLargeError(&TooLargeError{}))
	s.False(IsTooLargeError(&UnspecifiedTypeError{}))
	s.False(IsTooLargeError(errors.New("test")))
}

func (s *StorageErrorsSuite) TestIsUnspecifiedTypeError() {
	s.False(IsUnspecifiedTypeError(&UnregisteredTypeError{}))
	s.False(IsUnspecifiedTypeError(&InvalidPathError{}))
	s.False(IsUnspecifiedTypeError(&PathDoesntExistError{}))
	s.False(IsUnspecifiedTypeError(&TooLargeError{}))
	s.True(IsUnspecifiedTypeError(&UnspecifiedTypeError{}))
	s.False(IsUnspecifiedTypeError(errors.New("test")))
}

//
// Test Suite for the Type
//
type TypeSuite struct {
	suite.Suite
	storageType Type
}

func (s *TypeSuite) SetupSuite() {
	s.storageType = "MyTestingType"
	RegisterType(s.storageType, nil)
}

func (s *TypeSuite) TestTypeFmtString() {
	s.Equal(string(s.storageType), fmt.Sprintf("%s", s.storageType))
}

func (s *TypeSuite) TestTypeFmtV() {
	s.Equal(string(s.storageType), fmt.Sprintf("%v", s.storageType))
}

func (s *TypeSuite) TestRegisterTypeDuplicate() {
	s.Panics(func() {
		RegisterType(s.storageType, nil)
	})
}

func (s *TypeSuite) TestRegisterTypeUnspecified() {
	s.Panics(func() {
		RegisterType(Type(""), nil)
	})
}

func (s *TypeSuite) TestRegisterTypeTooLong() {
	tooLongType := Type(strings.Repeat("a", MaxTypeLen+1))
	s.Panics(func() {
		RegisterType(tooLongType, nil)
	})
}

//
// Test Suite for the New() function
//
type NewSuite struct {
	suite.Suite
}

func (s *NewSuite) TestUnspecifiedType() {
	c := &Conf{}
	st, err := New(c)
	s.Nil(st)
	s.True(IsUnspecifiedTypeError(err))
}

func (s *NewSuite) TestUnregisteredType() {
	c := &Conf{Type: Type("Doesn't exist")}
	st, err := New(c)
	s.Nil(st)
	s.NotNil(err)
	s.IsType(&UnregisteredTypeError{}, err)
}

func (s *NewSuite) TestNew() {
	myTestType := Type("TypeTestNew")
	factCalled := false
	fact := func(conf *Conf) (Storage, error) {
		factCalled = true
		return nil, nil
	}
	RegisterType(myTestType, fact)

	_, err := New(&Conf{Type: myTestType})
	s.Nil(err)
	s.True(factCalled)
}

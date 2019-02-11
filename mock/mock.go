package mock

import (
	"github.com/pw1/stor"
	"github.com/stretchr/testify/mock"
)

const (
	// MockStorageType is storage.Type of StorageMock.
	MockStorageType stor.Type = "Mock"
)

func init() {
	newStorageFunc := func(conf *stor.Conf) (stor.Storage, error) {
		return New(conf)
	}
	stor.RegisterType(MockStorageType, newStorageFunc)
}

// Mock is a mock object for mocking storage in testing.
type Mock struct {
	mock.Mock
}

// New creates a new storage.Mock
func New(conf *stor.Conf) (*Mock, error) {
	s := &Mock{}
	return s, nil
}

// Meta returns all entries within a directory.
func (m *Mock) Meta(path string) (*stor.Meta, error) {
	args := m.Called(path)
	return args.Get(0).(*stor.Meta), args.Error(1)
}

// List returns all entries within a directory.
func (m *Mock) List(path string) ([]string, []string, error) {
	args := m.Called(path)
	return args.Get(0).([]string), args.Get(1).([]string), args.Error(2)
}

// Load a file and return its content.
func (m *Mock) Load(path string, maxSize int64) ([]byte, error) {
	args := m.Called(path, maxSize)
	return args.Get(0).([]byte), args.Error(1)
}

// Save data to a file.
func (m *Mock) Save(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

// Delete a file.
func (m *Mock) Delete(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// Type returns the storage.Type of this storega object.
func (m *Mock) Type() stor.Type {
	args := m.Called()
	return args.Get(0).(stor.Type)
}

package mock

import (
	"testing"

	"github.com/pw1/stor"
	"github.com/stretchr/testify/suite"
)

// TestStorMockSuite is the test function that runs the tests in the StorMockSuite.
func TestStorMockSuite(t *testing.T) {
	suite.Run(t, new(StorMockSuite))
}

// TaskStorageSuite is the test suite for the StorMockSuite object.
type StorMockSuite struct {
	suite.Suite
}

// TestStorMockAsStorage makes sure that Mock actually implements the stor.Storage interface. If
// a method is missing or incorrect, then this won't compile.
func (s *StorMockSuite) TestStorMockAsStorage() {
	var storage stor.Storage
	storage, err := New(nil)
	s.NotNil(storage)
	s.Nil(err)
}

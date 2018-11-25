package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pw1/stor"
	"github.com/pw1/stor/tester"
)

// Call the generic storage tests
func TestMemoryStorageTester(t *testing.T) {
	testSuite := tester.New(MemoryStorageType)

	testSuite.SetupTestFunc = func(s *tester.StorageTester) stor.Storage {
		return New()
	}

	suite.Run(t, testSuite)
}

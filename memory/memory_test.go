package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pw1/stor"
	"github.com/pw1/stor/tester"
)

// TestMemoryStorageTester calls the generic storage tests
func TestMemoryStorageTester(t *testing.T) {
	myConfFactory := func() *stor.Conf {
		return &stor.Conf{
			Type: MemoryStorageType,
		}
	}

	testSuite := &tester.StorageTester{
		ConfFactory: myConfFactory,
	}

	suite.Run(t, testSuite)
}

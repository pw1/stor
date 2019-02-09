package amazons3

import (
	"errors"

	"github.com/pw1/stor"
)

const (
	// S3StorageType is the type of the S3 storage.
	S3StorageType stor.Type = "S3"
)

func init() {
	newStorageFunc := func(conf *stor.Conf) (stor.Storage, error) {
		return New(conf)
	}
	stor.RegisterType(S3StorageType, newStorageFunc)
}

// S3 is in implementation of stor.Storage. It uses Amazon's S3, or another compatible service, as
// it storage backend.
type S3 struct{}

// New create a new S3 object with the specified configuration.
func New(conf *stor.Conf) (*S3, error) {
	am := &S3{}
	return am, nil
}

// Meta returns meta information about a file.
func (s *S3) Meta(filePath string) (*stor.Meta, error) {
	return nil, errors.New("not yet implemented")
}

// List returns the files and subdirectories within the specified directory.
func (s *S3) List(path string) ([]string, []string, error) {
	return []string{}, []string{}, errors.New("not yet implemented")
}

// Load loads the content of the specified file. If the file is larger than maxSize, the an error is
// returned.
func (s *S3) Load(path string, maxSize int64) ([]byte, error) {
	return []byte{}, errors.New("not yet implemented")
}

// Save saves the data to the specified file.
func (s *S3) Save(path string, data []byte) error {
	return errors.New("not yet implemented")
}

// Delete removes a file from storage.
func (s *S3) Delete(path string) error {
	return errors.New("not yet implemented")
}

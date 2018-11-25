package amazons3

import (
	"fmt"

	"github.com/pw1/stor"
)

const (
	AmazonS3StorageType     stor.Type = 9149473504714319481
	AmazonS3StorageTypeText string    = "AmazonS3"
)

func init() {
	stor.RegisterStorageType(AmazonS3StorageType, AmazonS3StorageTypeText)
}

type AmazonS3 struct{}

func New() (*AmazonS3, error) {
	am := &AmazonS3{}
	return am, nil
}

func (a *AmazonS3) List(path string) ([]string, []string, error) {
	return []string{}, []string{}, fmt.Errorf("not yet implemented")
}

func (a *AmazonS3) Load(path string, maxSize int64) ([]byte, error) {
	return []byte{}, fmt.Errorf("not yet implemented")
}

func (a *AmazonS3) Save(path string, data []byte) error {
	return fmt.Errorf("not yet implemented")
}

func (a *AmazonS3) Type() stor.Type {
	return AmazonS3StorageType
}

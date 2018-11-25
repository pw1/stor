package stor

import (
	"errors"
	"fmt"
	"strings"
)

// Storage defines a simple interface for accessing different kinds of storage.
// The storage interface is for loading and saving blobs of data. The data is accessed via a
// hierarichal path. The directories within the path are separated by the slash '/' (even on Windows
// platforms).
type Storage interface {
	// List all entries within a directory.
	// The path argument is a slash-separated path.
	// Returns three values. The first return value is a list of files within the directory. The
	// second return value is the list of subdirectories within the directory. And the third return
	// value is any error that occured. The returned file and subdirectory entries are not
	// necessarily sorted. The returned file and subdirectory entries are always full paths (with
	// respect to the storage root).
	List(path string) ([]string, []string, error)

	// Check if a file exists.
	// The path argument is a slash-separated path.
	Exist(path string) (bool, error)

	// Load a file and return its content.
	// The path argument is a slash-separated path.
	// The maxSize gives the maximum accepted file size. If the file is larger, then an error is
	// returned and no data.
	Load(path string, maxSize int64) ([]byte, error)

	// Save data to a file.
	// The path argument is a slash-separated path.
	Save(path string, data []byte) error

	// Delete a file.
	// The path argument is a slash-separated path.
	Delete(path string) error

	// Return the StorageType of this storega object.
	Type() Type
}

// Type defines the type of Storage. Each type of storage has its own type ID. The Type has a
// numerical and a textual representation.
type Type int64

const (
	// UndefinedStorageType is the numerical representation of a storage.Type that is undefined.
	UndefinedStorageType Type = 0

	// UndefinedStorageTypeText is the textual representation of a storage.Type that is undefined.
	UndefinedStorageTypeText string = "Undefined"

	// MockStorageType is the storage.Type for the StorageMock.
	MockStorageType Type = 69844752684693179

	// MockStorageTypeText is the textual representation of the storage.Type of StorageMock.
	MockStorageTypeText string = "Mock"
)

var (
	// Defines the mapping between StorageTypes and their textual representations.
	registeredStorageTypes = make(map[string]Type)
)

func init() {
	RegisterStorageType(UndefinedStorageType, UndefinedStorageTypeText)
	RegisterStorageType(MockStorageType, MockStorageTypeText)
}

// RegisterStorageType registers a new storage.Type. You need to register a new storage type before
// it can be unmarshalled from text or formatted as string.
func RegisterStorageType(storageType Type, stringRepresentation string) error {
	// Make sure that int and text represenations are not already registered.
	for key, value := range registeredStorageTypes {
		if stringRepresentation == key {
			return fmt.Errorf("%s is already registered", stringRepresentation)
		}
		if int(value) == int(storageType) {
			return fmt.Errorf("%d is already registered (with %s)", storageType, value)
		}
	}

	registeredStorageTypes[stringRepresentation] = storageType
	return nil
}

// UnmarshalText parses a textual representation of a storage.Type and sets this object to that
// value.
func (s *Type) UnmarshalText(text []byte) error {
	typ, ok := registeredStorageTypes[string(text)]
	if !ok {
		msg := fmt.Sprintf("Invalid StorageType: %s (valid types are:", text)
		for typeName := range registeredStorageTypes {
			msg += fmt.Sprintf(" %s", typeName)
		}
		msg += ")."
		return errors.New(msg)
	}

	*s = typ
	return nil
}

func (s Type) String() string {
	for typeName, typ := range registeredStorageTypes {
		if s == typ {
			return typeName
		}
	}
	return fmt.Sprintf("INVALID TYPE: %d", s)
}

// Conf contains the configuration for the storege objects.
type Conf struct {
	StorageType Type
	Path        string
}

// NewConf creates a new empty configuration, with default values.
func NewConf() *Conf {
	return &Conf{
		StorageType: UndefinedStorageType,
		Path:        "",
	}
}

// InvalidPathError indicates that a path is invalid.
type InvalidPathError struct {
	msg string
}

// NewInvalidPathError generates a new InvalidPathError
func NewInvalidPathError(msg string) error {
	err := &InvalidPathError{
		msg: msg,
	}
	return err
}

// IsInvalidPathError checks whether an error is an InvalidPathError, or not.
func IsInvalidPathError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "InvalidPathError:")
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("InvalidPathError: %s", e.msg)
}

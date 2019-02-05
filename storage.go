package stor

import (
	"fmt"
)

// Metaer (Meta-er) can retrieve meta information about a file.
type Metaer interface {
	// Meta returns meta information about a file.
	// If the file does not exist, then a PathDoesntExistError is returned.
	Meta(path string) (*Meta, error)
}

// Lister can list files in Storage in a directory.
type Lister interface {
	// List all entries within a directory.
	// The path argument is a slash-separated path.
	// Returns three values. The first return value is a list of files within the directory. The
	// second return value is the list of subdirectories within the directory. And the third return
	// value is any error that occured. The returned file and subdirectory entries are not
	// necessarily sorted. The returned file and subdirectory entries are always full paths (with
	// respect to the storage root).
	List(path string) ([]string, []string, error)
}

// Loader can load files in Storage.
type Loader interface {
	// Load a file and return its content.
	// The path argument is a slash-separated path.
	// The maxSize gives the maximum accepted file size. If the file is larger, then an error is
	// returned and no data.
	Load(path string, maxSize int64) ([]byte, error)
}

// Saver can save files in Storage.
type Saver interface {
	// Save data to a file.
	// The path argument is a slash-separated path.
	Save(path string, data []byte) error
}

// Deleter can delete files from Storage.
type Deleter interface {
	// Delete a file.
	// The path argument is a slash-separated path.
	Delete(path string) error
}

// Reader can perform all read operations in Storage.
type Reader interface {
	Metaer
	Lister
	Loader
}

// Writer can perform all write operations in Storage.
type Writer interface {
	Saver
	Deleter
}

// Storage defines a simple, limited interface for accessing different kinds of storage.
// The storage interface is for loading and saving blobs of data. The data is accessed via a
// hierarichal path. The directories within the path are separated by the slash '/' (even on Windows
// platforms).
type Storage interface {
	Reader
	Writer
}

// Meta contains meta information about a file.
type Meta struct {
	// Size (in bytes) of the file. This value is set to SizeUnknown if the Size can't be retrieved.
	Size int64
}

const (
	// SizeUnknown indicates that the size of a file is unknown.
	SizeUnknown = -1
)

// Factory is a function that creates a new Storage object based on a configuration.
type Factory func(conf *Conf) (Storage, error)

// Type defines the type of Storage.
type Type string

const (
	// MaxTypeLen is the maximum length that a Type can have.
	MaxTypeLen = 20

	// TypeUnspecified indicates that the storage.Type is not specified.
	TypeUnspecified Type = ""
)

var (
	// typeFactoryMap contains the mapping between Types and their Factory functions.
	typeFactoryMap = make(map[Type]Factory)
)

// RegisterType registers a new storage.Type and its associated Factory function.
// If the Type is already registered, or if the Type is invalid, then this function will panic.
// This function is intended to be called from the init function of packages that implement the
// Storage interface.
func RegisterType(storageType Type, factory Factory) {
	if len(storageType) > MaxTypeLen {
		panic(fmt.Sprintf("stor: name of Type %s is too long", storageType))
	}

	if storageType == TypeUnspecified {
		panic("stor: undefined Type")
	}

	if _, ok := typeFactoryMap[storageType]; ok {
		panic(fmt.Sprintf("stor: Type %s is already registered", storageType))
	}

	typeFactoryMap[storageType] = factory
}

// New creates a new Storage object based on conf. It will read the Type from the conf and get the
// Factory function registered for that type. It will then call that Factory with conf and return
// the result.
func New(conf *Conf) (Storage, error) {
	if conf.Type == TypeUnspecified {
		return nil, &UnspecifiedTypeError{}
	}

	factory, ok := typeFactoryMap[conf.Type]
	if !ok {
		return nil, &UnregisteredTypeError{conf.Type}
	}

	return factory(conf)
}

// Conf contains the configuration for the storege objects.
type Conf struct {
	Type Type
	Path string
}

// UnregisteredTypeError is returned when a storage Type is specified but has never been registered.
type UnregisteredTypeError struct {
	Type Type
}

func (e *UnregisteredTypeError) Error() string {
	return fmt.Sprintf("storage type %s is not registered", e.Type)
}

// IsUnregisteredTypeError returns true if an error is a UnspecifiedTypeError. Returns false
// otherwise.
func IsUnregisteredTypeError(err error) bool {
	switch err.(type) {
	case *UnregisteredTypeError:
		return true
	default:
		return false
	}
}

// InvalidPathError indicates that a path is invalid.
type InvalidPathError struct {
	Path string
	Msg  string
}

func (e *InvalidPathError) Error() string {
	msg := fmt.Sprintf("path %s is invalid", e.Path)
	if e.Msg != "" {
		msg += ": " + e.Msg
	}
	return msg
}

// IsInvalidPathError checks whether an error is an InvalidPathError, or not.
func IsInvalidPathError(err error) bool {
	switch err.(type) {
	case *InvalidPathError:
		return true
	default:
		return false
	}
}

// PathDoesntExistError indicates that a specified path doesn't exist.
type PathDoesntExistError struct {
	// Path is the path that doesn't exist.
	Path string
}

func (f *PathDoesntExistError) Error() string {
	return fmt.Sprintf("path %s does not exist", f.Path)
}

// IsPathDoesntExistError returns true if an error is a PathDoesntExistError. Returns false
// otherwise.
func IsPathDoesntExistError(err error) bool {
	switch err.(type) {
	case *PathDoesntExistError:
		return true
	default:
		return false
	}
}

// TooLargeError indicates that a file is too large, or a list is too long.
type TooLargeError struct {
	// What indicates what is too large. E.g. a file or a list.
	What string
}

func (e *TooLargeError) Error() string {
	msg := "too large"
	if e.What != "" {
		msg = e.What + " is " + msg
	}
	return msg
}

// IsTooLargeError returns true if an error is a TooLargeError. Returns false otherwise.
func IsTooLargeError(err error) bool {
	switch err.(type) {
	case *TooLargeError:
		return true
	default:
		return false
	}
}

// UnspecifiedTypeError is returned when trying to create Storage but Type is not specified.
type UnspecifiedTypeError struct{}

func (e *UnspecifiedTypeError) Error() string {
	return "storage Type is not specified"
}

// IsUnspecifiedTypeError returns true if an error is a UnspecifiedTypeError. Returns false
// otherwise.
func IsUnspecifiedTypeError(err error) bool {
	switch err.(type) {
	case *UnspecifiedTypeError:
		return true
	default:
		return false
	}
}

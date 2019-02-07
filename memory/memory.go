// Package memory implements the stor.Storage interface as a map in memory.
package memory

import (
	"strings"

	"github.com/pw1/stor"
)

const (
	// MemoryStorageType is the storage type of the Memory storage.
	MemoryStorageType stor.Type = "Memory"
)

func init() {
	newStorageFunc := func(conf *stor.Conf) (stor.Storage, error) {
		return New(conf)
	}
	stor.RegisterType(MemoryStorageType, newStorageFunc)
}

// Memory is a stor.Storage implementation. It stores everything in memory. Can, for example, be
// used as memory cache, or for testing.
type Memory struct {
	data map[string][]byte
}

// New creates a new Memory storage.
// The supplied configuration has not effect on the created Memory object.
func New(conf *stor.Conf) (*Memory, error) {
	mem := &Memory{
		data: make(map[string][]byte),
	}
	return mem, nil
}

// Meta returns meta information about a file.
func (m *Memory) Meta(filePath string) (*stor.Meta, error) {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return nil, err
	}

	data, ok := m.data[cleanPath]
	if !ok {
		return nil, &stor.PathDoesntExistError{Path: cleanPath}
	}

	meta := &stor.Meta{
		Size: int64(len(data)),
	}

	return meta, nil
}

// List returns the files and subdirectories within the specified directory.
func (m *Memory) List(filePath string) ([]string, []string, error) {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return []string{}, []string{}, err
	}

	var prefix string
	if cleanPath == "" {
		prefix = cleanPath
	} else {
		prefix = cleanPath + "/"
	}

	files := make([]string, 0)
	dirsMap := make(map[string]bool)
	for key := range m.data {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		withoutPrefix := key[len(prefix):]
		slashIdx := strings.Index(withoutPrefix, "/")
		if slashIdx < 0 {
			fullPath := prefix + withoutPrefix
			files = append(files, fullPath)
		} else {
			fullPath := prefix + withoutPrefix[:slashIdx]
			dirsMap[fullPath] = true
		}
	}

	// Convert the map with directories to a slice. We used the map to avoid duplicates
	dirs := make([]string, 0, len(dirsMap))
	for dir := range dirsMap {
		dirs = append(dirs, dir)
	}

	return files, dirs, nil
}

// Load loads the content of the specified file. If the file is larger than maxSize, the an error is
// returned.
func (m *Memory) Load(filePath string, maxSize int64) ([]byte, error) {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return []byte{}, err
	}

	dataInStorage, ok := m.data[cleanPath]
	if !ok {
		return []byte{}, &stor.PathDoesntExistError{Path: cleanPath}
	}

	if int64(len(dataInStorage)) > maxSize {
		return []byte{}, &stor.TooLargeError{What: cleanPath}
	}

	dataCopy := make([]byte, len(dataInStorage))
	copy(dataCopy, dataInStorage)

	return dataCopy, nil
}

// Save saves the data to the specified file.
func (m *Memory) Save(filePath string, data []byte) error {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return err
	}

	m.data[cleanPath] = make([]byte, len(data))
	copy(m.data[cleanPath], data)

	return nil
}

// Delete removes a file from storage.
func (m *Memory) Delete(filePath string) error {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return err
	}

	if _, ok := m.data[cleanPath]; !ok {
		return &stor.PathDoesntExistError{
			Path: cleanPath,
		}
	}

	delete(m.data, cleanPath)
	return nil
}

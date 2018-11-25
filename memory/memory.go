package memory

import (
	"fmt"
	"strings"

	"github.com/pw1/stor"
)

const (
	MemoryStorageType     stor.Type = 3331380253917538617
	MemoryStorageTypeText string    = "Memory"
)

func init() {
	stor.RegisterStorageType(MemoryStorageType, MemoryStorageTypeText)
}

type Memory struct {
	data map[string][]byte
}

func New() *Memory {
	mem := &Memory{
		data: make(map[string][]byte),
	}
	return mem
}

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

func (m *Memory) Exist(filePath string) (bool, error) {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return false, err
	}

	_, ok := m.data[cleanPath]

	return ok, nil
}

func (m *Memory) Load(filePath string, maxSize int64) ([]byte, error) {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return []byte{}, err
	}

	dataInStorage, ok := m.data[cleanPath]
	if !ok {
		return []byte{}, fmt.Errorf("%s does not exist", filePath)
	}

	if int64(len(dataInStorage)) > maxSize {
		return []byte{}, fmt.Errorf("Data too large. Data is %d bytes, while maxSize is %d",
			len(dataInStorage), maxSize)
	}

	dataCopy := make([]byte, len(dataInStorage))
	copy(dataCopy, dataInStorage)

	return dataCopy, nil
}

func (m *Memory) Save(filePath string, data []byte) error {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return err
	}

	m.data[cleanPath] = make([]byte, len(data))
	copy(m.data[cleanPath], data)

	return nil
}

func (m *Memory) Delete(filePath string) error {
	cleanPath, err := stor.CleanPath(filePath)
	if err != nil {
		return err
	}

	delete(m.data, cleanPath)
	return nil
}

func (m *Memory) Type() stor.Type {
	return MemoryStorageType
}

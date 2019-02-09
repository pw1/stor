package localdir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pw1/stor"
)

const (
	// LocalDirStorageType is the storage type of the LocalDir storage.
	LocalDirStorageType stor.Type = "LocalDir"
)

func init() {
	newStorageFunc := func(conf *stor.Conf) (stor.Storage, error) {
		return New(conf)
	}
	stor.RegisterType(LocalDirStorageType, newStorageFunc)
}

// LocalDir is a Storage object that uses a directory in the local file system as storage backend.
type LocalDir struct {
	BaseDir string
}

// New creates a new LocalDir object.
func New(conf *stor.Conf) (*LocalDir, error) {
	absPath, err := filepath.Abs(conf.Path)
	if err != nil {
		return nil, fmt.Errorf("Invalid base dir %v: %v", conf.Path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to use local dir %v: %v", absPath, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("Local dir %v is not a directory", absPath)
	}

	ldir := &LocalDir{
		BaseDir: absPath,
	}

	return ldir, nil
}

// Get the full absolute path. this function also checks whether the path escapes the BaseDir. An
// error is raised if the path escapes the BaseDir.
// The filePath argument is the slash-separated path.
// The returned path uses the platform specific directory separator '/' or '\'.
func (l *LocalDir) getFullPath(filePath string) (string, error) {
	filePath, err := stor.CleanPath(filePath)
	if err != nil {
		return "", err
	}

	// Convert the slash-separated path to an absolute, platform-dependent path
	fullPath, err := filepath.Abs(filepath.Join(l.BaseDir, filepath.FromSlash(filePath)))
	if err != nil {
		return "", fmt.Errorf("invalid filePath %v: %v", filePath, err)
	}

	// Double-check that we don't escape from the base directory
	if escapesDir(fullPath, l.BaseDir) {
		msg := fmt.Sprintf("invalid filePath %v, it escapes the base directory", filePath)
		return "", &stor.InvalidPathError{Path: msg}
	}

	return fullPath, nil
}

// Meta returns meta information about a file.
func (l *LocalDir) Meta(filePath string) (*stor.Meta, error) {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &stor.PathDoesntExistError{Path: filePath}
		}
		return nil, err
	}

	meta := &stor.Meta{
		Size: info.Size(),
	}

	return meta, nil
}

// List returns the files and subdirectories within the specified directory.
func (l *LocalDir) List(filePath string) ([]string, []string, error) {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return []string{}, []string{}, err
	}

	entries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return []string{}, []string{}, err
	}

	files := []string{}
	dirs := []string{}
	for _, entry := range entries {
		slashPathWithinStorage := path.Join(filePath, entry.Name())
		if entry.IsDir() {
			dirs = append(dirs, slashPathWithinStorage)
		} else {
			files = append(files, slashPathWithinStorage)
		}
	}

	return files, dirs, nil
}

// Load loads the content of the specified file. If the file is larger than maxSize, the an error is
// returned.
func (l *LocalDir) Load(filePath string, maxSize int64) ([]byte, error) {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return []byte{}, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, &stor.PathDoesntExistError{Path: filePath}
		}
		return []byte{}, err
	}

	if info.Size() > maxSize {
		return []byte{}, &stor.TooLargeError{What: filePath}
	}

	return ioutil.ReadFile(fullPath)
}

// Save saves the data to the specified file.
func (l *LocalDir) Save(filePath string, data []byte) error {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return err
	}

	// Make sure that the parent directory exists
	dirPath := filepath.Dir(fullPath)
	err = os.MkdirAll(dirPath, 0700)

	err = ioutil.WriteFile(fullPath, data, 0660)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a file from storage.
func (l *LocalDir) Delete(filePath string) error {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return err
	}

	err = os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &stor.PathDoesntExistError{Path: filePath}
		}
		return err
	}

	// Remove all empty parent directories (until we reach the basedir)
	parentDir := fullPath
	for i := 0; true; i++ {
		if i > 1000 {
			return fmt.Errorf("Infinite loop in LocalDir.Delete()")
		}

		parentDir = filepath.Dir(parentDir)
		if escapesDir(parentDir, l.BaseDir) || (parentDir == l.BaseDir) {
			break
		}

		entries, err := ioutil.ReadDir(parentDir)
		if err != nil {
			return err
		}
		if len(entries) > 0 {
			break
		}
		os.Remove(parentDir)
	}

	return nil
}

// escapesDir checks whether a path escapes a certain baseDir directory.
// Return true if path is not within the baseDir. Returns false if path is within the baseDir, or
// equal to baseDir.
func escapesDir(path, baseDir string) bool {
	path, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return true
	}

	baseDirSlash := addTrailingSeparator(baseDir)
	pathSlash := addTrailingSeparator(path)

	return !strings.HasPrefix(pathSlash, baseDirSlash)
}

// addTrailingSeparator add a trailing separator (e.g. / on Linux and \ on Windows) to a path if
// that path does not yet have such separator at the end.
func addTrailingSeparator(path string) string {
	n := len(path)
	if (n == 0) || path[n-1] != filepath.Separator {
		path = path + string(filepath.Separator)
	}

	return path
}

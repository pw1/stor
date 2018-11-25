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
	// LocalDirStorageType defines the numeric representation of the LocalDir stor.Type.
	LocalDirStorageType stor.Type = 3707851827220653854

	// LocalDirStorageTypeText defines the textual representation of the LocalDir stor.Type.
	LocalDirStorageTypeText string = "LocalDir"
)

func init() {
	stor.RegisterStorageType(LocalDirStorageType, LocalDirStorageTypeText)
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

	if ldir.Type() != conf.StorageType {
		return nil, fmt.Errorf("invalid StorageType %s. I'm expecting type %s", conf.StorageType,
			ldir.Type())
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
		return "", stor.NewInvalidPathError(msg)
	}

	return fullPath, nil
}

// List all entries within a directory.
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

// Exist checks whether a file exists, or not.
func (l *LocalDir) Exist(filePath string) (bool, error) {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

// Load the content of a file. Return an error if the file is larger than maxSize.
func (l *LocalDir) Load(filePath string, maxSize int64) ([]byte, error) {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return []byte{}, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return []byte{}, err
	}

	if info.Size() > maxSize {
		return []byte{}, fmt.Errorf("File is larger than %d", maxSize)
	}

	return ioutil.ReadFile(fullPath)
}

// Save the content of a file.
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

// Delete a file from stor.
func (l *LocalDir) Delete(filePath string) error {
	fullPath, err := l.getFullPath(filePath)
	if err != nil {
		return err
	}

	err = os.Remove(fullPath)
	if err != nil {
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

// Type returns the Type of this Storage object.
func (l *LocalDir) Type() stor.Type {
	return LocalDirStorageType
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

func addTrailingSeparator(path string) string {
	n := len(path)
	if (n == 0) || path[n-1] != filepath.Separator {
		path = path + string(filepath.Separator)
	}

	return path
}

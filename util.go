package stor

import (
	"fmt"
	"path"
	"strings"
)

const (
	// Delimiter separates path components
	Delimiter = '/'

	// ValidBytes lists the bytes (characters) that are valid in path components.
	ValidBytes = "._-abcdefghijklmnopqrstuwvxyzABCDEFGHIJKLMNOPQRSTUWVXYZ0123456789"
)

var (
	// Forbidden combinations
	Forbidden = []string{
		"..",
	}

	validCharDict = make(map[byte]bool)
)

func init() {
	// Create a dictionary with all allowed characters. This is to allow quick lookup.
	validCharDict[Delimiter] = true
	for i := 0; i < len(ValidBytes); i++ {
		validCharDict[ValidBytes[i]] = true
	}
}

// CleanPath cleans up a path for use in Storage objects.
func CleanPath(filePath string) (string, error) {
	// Check for any forbidden combinations
	for _, forbid := range Forbidden {
		if strings.Contains(filePath, forbid) {
			msg := fmt.Sprintf("Path contains forbidden combination %s", forbid)
			return "", NewInvalidPathError(msg)
		}
	}

	// Make sure that the path doesn't start with a /
	if path.IsAbs(filePath) {
		msg := fmt.Sprintf("Absolute path %s", filePath)
		return "", NewInvalidPathError(msg)
	}

	// Check for any forbidden characters
	for i := 0; i < len(filePath); i++ {
		char := filePath[i]
		_, ok := validCharDict[char]
		if !ok {
			msg := fmt.Sprintf("%s contains forbidden byte 0x%x (%s) at index %d", filePath,
				char, string(char), i)
			return "", NewInvalidPathError(msg)
		}
	}

	// Clean the path (removing any // combinations)
	cleanPath := path.Clean(filePath)
	if cleanPath == "." {
		cleanPath = ""
	}

	return cleanPath, nil
}

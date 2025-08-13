package server

import (
	"errors"
	"path/filepath"
	"strings"
)

// resolveWithinRoot cleans and joins a user-supplied relative path to root
// and ensures the final path remains within root to prevent directory traversal.
func resolveWithinRoot(root string, rel string) (string, error) {
    cleaned := filepath.Clean("/" + rel)
    cleaned = strings.TrimPrefix(cleaned, "/")
    abs := filepath.Join(root, cleaned)
    // Ensure the resulting path stays under root
    rootAbs, _ := filepath.Abs(root)
    absClean, _ := filepath.Abs(abs)
    if !strings.HasPrefix(absClean+string(filepath.Separator), rootAbs+string(filepath.Separator)) && absClean != rootAbs {
        return "", errors.New("path escapes root")
    }
    return absClean, nil
}



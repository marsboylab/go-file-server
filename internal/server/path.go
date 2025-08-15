package server

import (
	"errors"
	"path/filepath"
	"strings"
)

func resolveWithinRoot(root string, rel string) (string, error) {
    cleaned := filepath.Clean("/" + rel)
    cleaned = strings.TrimPrefix(cleaned, "/")
    abs := filepath.Join(root, cleaned)
    rootAbs, _ := filepath.Abs(root)
    absClean, _ := filepath.Abs(abs)
    if !strings.HasPrefix(absClean+string(filepath.Separator), rootAbs+string(filepath.Separator)) && absClean != rootAbs {
        return "", errors.New("path escapes root")
    }
    return absClean, nil
}



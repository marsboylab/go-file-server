package server

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
    rootDir string
}

func NewFileService(root string) *FileService {
    return &FileService{rootDir: root}
}

type FileInfo struct {
    Name         string    `json:"name"`
    Size         int64     `json:"size"`
    ModTime      time.Time `json:"modTime"`
    RelativePath string    `json:"relativePath"`
    IsDir        bool      `json:"isDir"`
}

func (s *FileService) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
    target := r.URL.Query().Get("path")
    if target == "" {
        target = "."
    }
    dir, err := resolveWithinRoot(s.rootDir, target)
    if err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }

    entries, err := os.ReadDir(dir)
    if err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }

    var files []FileInfo
    for _, e := range entries {
        info, err := e.Info()
        if err != nil {
            continue
        }
        files = append(files, FileInfo{
            Name:         e.Name(),
            Size:         info.Size(),
            ModTime:      info.ModTime(),
            RelativePath: filepath.ToSlash(filepath.Join(target, e.Name())),
            IsDir:        e.IsDir(),
        })
    }

    sort.Slice(files, func(i, j int) bool { return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name) })
    writeJSON(w, http.StatusOK, files)
}

func (s *FileService) UploadHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(64 << 20); err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }

    target := r.FormValue("path")
    if target == "" {
        target = "."
    }
    dir, err := resolveWithinRoot(s.rootDir, target)
    if err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }
    if err := os.MkdirAll(dir, 0o755); err != nil {
        writeError(w, http.StatusInternalServerError, err)
        return
    }

    files := r.MultipartForm.File["file"]
    if len(files) == 0 {
        writeError(w, http.StatusBadRequest, errors.New("file field is required"))
        return
    }

    var saved []FileInfo
    for _, fh := range files {
        f, err := fh.Open()
        if err != nil {
            writeError(w, http.StatusBadRequest, err)
            return
        }
        defer f.Close()

        filename := sanitizeFilename(fh.Filename)
        if filename == "" {
            filename = uuid.NewString()
        }
        dstPath := filepath.Join(dir, filename)

        out, err := os.Create(dstPath)
        if err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }
        if _, err := io.Copy(out, f); err != nil {
            out.Close()
            writeError(w, http.StatusInternalServerError, err)
            return
        }
        if err := out.Close(); err != nil {
            writeError(w, http.StatusInternalServerError, err)
            return
        }

        fi, _ := os.Stat(dstPath)
        saved = append(saved, FileInfo{
            Name:         filename,
            Size:         fi.Size(),
            ModTime:      fi.ModTime(),
            RelativePath: filepath.ToSlash(filepath.Join(target, filename)),
            IsDir:        false,
        })
    }

    writeJSON(w, http.StatusCreated, saved)
}

func (s *FileService) DownloadHandler(w http.ResponseWriter, r *http.Request) {
    rel := strings.TrimPrefix(r.URL.Path, "/api/files/")
    if rel == "" || rel == "/" || rel == "*" {
        writeError(w, http.StatusBadRequest, errors.New("invalid file path"))
        return
    }
    target, err := resolveWithinRoot(s.rootDir, rel)
    if err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }

    fi, err := os.Stat(target)
    if err != nil || fi.IsDir() {
        writeError(w, http.StatusNotFound, errors.New("file not found"))
        return
    }

    http.ServeFile(w, r, target)
}

func (s *FileService) DeleteHandler(w http.ResponseWriter, r *http.Request) {
    rel := strings.TrimPrefix(r.URL.Path, "/api/files/")
    if rel == "" || rel == "/" || rel == "*" {
        writeError(w, http.StatusBadRequest, errors.New("invalid file path"))
        return
    }
    target, err := resolveWithinRoot(s.rootDir, rel)
    if err != nil {
        writeError(w, http.StatusBadRequest, err)
        return
    }
    if err := os.Remove(target); err != nil {
        writeError(w, http.StatusNotFound, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
    type resp struct {
        Error string `json:"error"`
    }
    writeJSON(w, status, resp{Error: err.Error()})
}

func sanitizeFilename(name string) string {
    name = filepath.Base(name)
    name = strings.TrimSpace(name)
    name = strings.ReplaceAll(name, "..", "")
    name = strings.ReplaceAll(name, string(filepath.Separator), "-")
    return name
}

// ensure multipart import is used in the build
var _ = multipart.FileHeader{}



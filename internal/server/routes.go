package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, fs *FileService) {
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    r.Get("/openapi.json", OpenAPIHandler)
    r.Get("/docs", SwaggerUIHandler)

    r.Route("/api", func(r chi.Router) {
        r.Get("/files", fs.ListFilesHandler)
        r.Get("/files/*", fs.DownloadHandler)
        r.Post("/files", fs.UploadHandler)
        r.Delete("/files/*", fs.DeleteHandler)
    })
}



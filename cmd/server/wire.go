package main

import (
	"github.com/go-chi/chi/v5"
	srv "github.com/user/go-file-server/internal/server"
)

func NewFileService(root string) *srv.FileService { return srv.NewFileService(root) }
func registerRoutes(r *chi.Mux, fs *srv.FileService) { srv.RegisterRoutes(r, fs) }



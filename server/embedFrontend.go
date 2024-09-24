package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend
//go:embed frontend/_next
//go:embed frontend/_next/static/chunks/pages/*.js
//go:embed frontend/_next/static/*/*.js
var frontendFiles embed.FS

func getFrontendFileSystem() http.FileSystem {
	fsys, err := fs.Sub(frontendFiles, "frontend")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

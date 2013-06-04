package gobioweb

import (
	"html/template"
	"os"
	"path/filepath"
)

type TemplateCache struct {
	Cache       *template.Template
	GlobPattern string
	Path        string
}

func NewTemplateCache() (tc *TemplateCache, err error) {
	tc = &TemplateCache{}
	tc.GlobPattern = "*.tmpl"
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	tc.Path = dir
	return
}

func (t *TemplateCache) CacheEntriesFromPath(path string) {
	fullPath := filepath.Join(path, t.GlobPattern)
	t.Cache = template.Must(template.ParseGlob(fullPath))
}

func (t *TemplateCache) CacheEntriesFromGlob(glob string ) {
	fullPath := filepath.Join(t.Path, glob)
	t.Cache = template.Must(template.ParseGlob(fullPath))
}

func (t *TemplateCache) CacheEntries() {
	fullPath := filepath.Join(t.Path, t.GlobPattern)
	t.Cache = template.Must(template.ParseGlob(fullPath))
}

package gobioweb

import (
	"html/template"
	"path/filepath"
)

var cachedEntries = map[string]*template.Template{}

type TemplateWithLayout struct {
	Layout    string
	Extension string
	Path      string
}

func (tl *TemplateWithLayout) Process(name string) *template.Template {
	if t, ok := cachedEntries[name]; ok {
		return t
	}
	tname := filepath.Join(tl.Path, name+tl.Extension)
	lname := filepath.Join(tl.Path, tl.Layout+tl.Extension)
	t := template.Must(template.New(tl.Layout).ParseFiles(lname, tname))
	cachedEntries[name] = t
	return t
}

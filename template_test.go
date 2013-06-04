package gobioweb

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

type templateFile struct {
	name     string
	contents string
}

func createTemplateFiles() []templateFile {
	files := []templateFile{
		{name: "temp1.tmpl",
			contents: `{{define "T1"}}
			             The confernce is in maldives for 
			             2 days with 50 people.
			           {{end}}
								 {{define "T2"}}
								 	 {{template "T1"}}
								 {{end}}
								 {{template "T2"}}
		`},
		{name: "temp2.tmpl",
			contents: `{{define "T3"}}
			             It is a event in keys for 
			             2 days with 50 people.
			          {{end}}
								{{define "T4"}}
								 	 {{template "T3"}}
								{{end}}
								{{template "T4"}}
		 `},
	}
	return files
}

func createTemplateFilesFromExt(ext string) []templateFile {
	files := []templateFile{
		{name: "temp1." + ext,
			contents: `{{define "T1"}}
			             The confernce is in maldives for 
			             2 days with 50 people.
			           {{end}}
								 {{define "T2"}}
								 	 {{template "T1"}}
								 {{end}}
								 {{template "T2"}}
		`},
		{name: "temp2." + ext,
			contents: `{{define "T3"}}
			             It is a event in keys for 
			             2 days with 50 people.
			          {{end}}
								{{define "T4"}}
								 	 {{template "T3"}}
								{{end}}
								{{template "T4"}}
		 `},
	}
	return files
}

func SetUpTemplateDir(files []templateFile) string {
	dir, err := ioutil.TempDir("", "template")
	if err != nil {
		log.Fatal(err)
	}
	for _, template := range files {
		f, err := os.Create(filepath.Join(dir, template.name))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = io.WriteString(f, template.contents)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir

}

func TestNewTemplateCache(t *testing.T) {
	r := regexp.MustCompile(`^\S+`)
	tc, _ := NewTemplateCache()
	if tc.GlobPattern != "*.tmpl" {
		t.Errorf("Expected *.tmpl, got %s", tc.GlobPattern)
	}
	if !r.MatchString(tc.Path) {
		t.Errorf("Expected path, got %s", tc.Path)
	}
}

func TestCacheEntries(t *testing.T) {
	files := createTemplateFiles()
	dir := SetUpTemplateDir(files)
	defer os.RemoveAll(dir)

	tc := &TemplateCache{Path: dir, GlobPattern: "*.tmpl"}
	tc.CacheEntries()

	var out bytes.Buffer
	err := tc.Cache.ExecuteTemplate(&out, "temp1.tmpl", nil)
	if err != nil {
		t.Error("Expected to execute template temp1.tmpl, did not happen")
	}

	err = tc.Cache.ExecuteTemplate(&out, "temp2.tmpl", nil)
	if err != nil {
		t.Error("Expected to execute template temp2.tmpl, did not happen")
	}
}

func TestCacheEntriesFromPath(t *testing.T) {
	files := createTemplateFiles()
	dir := SetUpTemplateDir(files)
	defer os.RemoveAll(dir)

	tc := &TemplateCache{GlobPattern: "*.tmpl"}
	tc.CacheEntriesFromPath(dir)

	var out bytes.Buffer
	err := tc.Cache.ExecuteTemplate(&out, "temp1.tmpl", nil)
	if err != nil {
		t.Error("Expected to execute template temp1.tmpl, did not happen")
	}

	err = tc.Cache.ExecuteTemplate(&out, "temp2.tmpl", nil)
	if err != nil {
		t.Error("Expected to execute template temp2.tmpl, did not happen")
	}
}

func TestCacheEntriesFromGlob(t *testing.T) {
	files := createTemplateFilesFromExt("html")
	dir := SetUpTemplateDir(files)
	defer os.RemoveAll(dir)

	tc := &TemplateCache{Path: dir}
	tc.CacheEntriesFromGlob("*.html")

	var out bytes.Buffer
	err := tc.Cache.ExecuteTemplate(&out, "temp1.html", nil)
	if err != nil {
		t.Error("Expected to execute template temp1.html, did not happen")
	}

	err = tc.Cache.ExecuteTemplate(&out, "temp2.html", nil)
	if err != nil {
		t.Error("Expected to execute template html.tmpl, did not happen")
	}
}

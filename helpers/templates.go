package helpers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"sync"
)

type TemplateCache struct {
	fs  fs.FS
	fm  template.FuncMap
	set map[string]*template.Template
	mu  sync.RWMutex
}

func NewTemplateCache(fsys fs.FS, fm template.FuncMap) *TemplateCache {
	return &TemplateCache{
		fs:  fsys,
		fm:  fm,
		set: make(map[string]*template.Template),
	}
}

func (tc *TemplateCache) Get(key, pageFile string) (*template.Template, error) {
	tc.mu.RLock()
	if t, ok := tc.set[key]; ok {
		tc.mu.RUnlock()
		return t, nil
	}
	tc.mu.RUnlock()

	base, err := template.New("layout").Funcs(tc.fm).ParseFS(
		tc.fs,
		"templates/layout.tmpl.html",
		"templates/partials/*.tmpl.html",
	)
	if err != nil {
		return nil, err
	}

	t, err := base.ParseFS(tc.fs, pageFile)
	if err != nil {
		return nil, err
	}

	tc.mu.Lock()
	tc.set[key] = t
	tc.mu.Unlock()
	return t, nil
}

func (tc *TemplateCache) Render(w http.ResponseWriter, name string, data any) error {
	pageFile := fmt.Sprintf("templates/%s.tmpl.html", name)
	t, err := tc.Get(name, pageFile)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, name, data)
}

package templates

import (
	"embed"
	"html/template"
	"io/fs"
	"strings"
	"time"
)

//go:embed *.html *.css
var content embed.FS

var pages map[string]*template.Template

var funcs = template.FuncMap{
	"date": func(t time.Time) string {
		return strings.ToUpper(t.Format("Jan 02"))
	},
	"time": func(t time.Time) string {
		return t.Format("3:04 PM")
	},
}

func init() {
	layout := template.Must(
		template.New("layout.html").Funcs(funcs).ParseFS(content, "layout.html", "layout.css"),
	)

	files, err := fs.Glob(content, "*.html")
	if err != nil {
		panic(err)
	}

	pages = make(map[string]*template.Template, len(files))
	for _, path := range files {
		name := strings.TrimSuffix(path, ".html")
		if name == "layout" {
			continue
		}

		t := template.Must(layout.Clone())
		template.Must(t.ParseFS(content, path))
		pages[name] = t
	}
}

// Get returns the parsed template for the given page name.
func Get(name string) *template.Template {
	return pages[name]
}

package models

import (
	"cmp"
	"io"
	"io/fs"
	"slices"
	"strings"
	"text/template"
	"time"

	"knowhere.cafe/src/shared"
)

const LAYOUT_PATH = "_layout.html"

type TemplateState map[string]*template.Template

func SetupTemplates(templateFiles fs.FS) TemplateState {
	state := make(TemplateState, 10)
	funcs := template.FuncMap{}
	ignored := []string{LAYOUT_PATH}

	fs.WalkDir(
		templateFiles,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || slices.Contains(ignored, path) {
				return nil
			}

			t, err := template.ParseFS(templateFiles, LAYOUT_PATH, path)
			if err != nil {
				return err
			}

			state[path] = t.Funcs(funcs)
			return nil
		},
	)

	return state
}

type TemplateData struct {
	Now      time.Time
	PageName string
	Auth     bool
	Data     any
}

func (ts TemplateState) Render(
	wr io.Writer,
	name string,
	target string,
	auth bool,
	data any,
) error {
	target = cmp.Or(target, LAYOUT_PATH)
	pageName := strings.TrimSuffix(name, ".html")

	td := TemplateData{
		time.Now(), pageName, auth, data,
	}

	t, ok := ts[name]
	if !ok {
		return shared.ErrUnknownTemplate{Name: name}
	}

	return t.ExecuteTemplate(wr, target, td)
}

type Renderable interface {
	Title() string
	Body() string
	PublishedAt() time.Time
	Markdown(w io.Writer) error
}

package models

import (
	"cmp"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"slices"
	"strings"
	"time"

	"knowhere.cafe/src/shared"
)

const LAYOUT_PATH = "_layout.html"
const MAIN_PATH = "_main.html"

var ignored = []string{LAYOUT_PATH, MAIN_PATH}

type TemplateState struct {
	inner map[string]*template.Template
	dev   bool
	// Doing something tricky here to let me only pass in the fs once
	// and then re-parse templates easily later in dev mode
	curried func(path string) error
}

var funcs = template.FuncMap{
	"safe": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func (ts TemplateState) setupTemplate(templateFiles fs.FS, path string) error {
	t, err := template.New("").
		Funcs(funcs).
		ParseFS(templateFiles, LAYOUT_PATH, MAIN_PATH, path)
	if err != nil {
		slog.Error("template parse error", "error", err)
		return err
	}

	ts.inner[path] = t
	return nil
}

func SetupTemplates(templateFiles fs.FS, dev bool) *TemplateState {
	state := TemplateState{
		inner: make(map[string]*template.Template, 10),
		dev:   dev,
	}

	state.curried = func(path string) error {
		return state.setupTemplate(templateFiles, path)
	}

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

			return state.curried(path)
		},
	)

	return &state
}

type TemplateData struct {
	Now      time.Time
	PageName string
	Auth     bool
	Data     any
}

func (ts TemplateState) Render(
	wr io.Writer,
	path string,
	target string,
	auth bool,
	data any,
) error {
	// re-parse the template in dev mode, which is when setupCurried exists
	if ts.dev {
		ts.curried(path)
	}

	target = cmp.Or(target, LAYOUT_PATH)
	pageName := strings.TrimSuffix(path, ".html")

	td := TemplateData{
		time.Now(), pageName, auth, data,
	}

	t, ok := ts.inner[path]
	if !ok {
		return shared.ErrUnknownTemplate{Name: path}
	}

	// This is a workaround because the `main` template DOES exist in _layout.html
	// and WILL be called by Unpoly but we actually don't want to use it
	// and specifically want to use _main.html
	if target == "main" {
		target = MAIN_PATH
	}

	return t.ExecuteTemplate(wr, target, td)
}

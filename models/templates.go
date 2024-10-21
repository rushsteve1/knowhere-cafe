package models

import (
	"cmp"
	"io"
	"io/fs"
	"log/slog"
	"slices"
	"strings"
	"html/template"
	"time"

	"knowhere.cafe/src/shared"
)

const LAYOUT_PATH = "_layout.html"

var ignored = []string{LAYOUT_PATH}

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
	t, err := template.ParseFS(templateFiles, LAYOUT_PATH, path)
	if err != nil {
		slog.Error("template parse error", "error", err)
		return err
	}

	ts.inner[path] = t.Funcs(funcs)
	return nil
}

func SetupTemplates(templateFiles fs.FS, dev bool) TemplateState {
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

	return t.ExecuteTemplate(wr, target, td)
}

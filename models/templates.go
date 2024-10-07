package models

import (
	"context"
	"io"
	"io/fs"
	"slices"
	"strings"
	"text/template"
	"time"

	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

type TemplateState = map[string]*template.Template

func SetupTemplates(templateFiles fs.FS) TemplateState {
	state := make(TemplateState, 10)
	funcs := template.FuncMap{}
	ignored := []string{shared.LAYOUT_PATH}

	fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || slices.Contains(ignored, path) {
			return nil
		}

		t, err := template.ParseFS(templateFiles, shared.LAYOUT_PATH, path)
		if err != nil {
			return err
		}

		state[path] = t.Funcs(funcs)
		return nil
	})

	return state
}

type TemplateData[T any] struct {
	Fragment bool
	Cfg      Config
	Now      time.Time
	PageName string
	Data     any
}

func Render[T any](
	ctx context.Context,
	wr io.Writer,
	name string,
	frag bool,
	data T,
) error {
	state := easy.Must(State(ctx))
	cfg, err := state.Config()
	if err != nil {
		return err
	}

	pageName := strings.TrimSuffix(name, ".html")

	td := TemplateData[T]{
		frag, cfg, time.Now(), pageName, data,
	}

	t, ok := state.Templ[name]
	if !ok {
		return shared.ErrUnknownTemplate{Name: name}
	}

	return t.ExecuteTemplate(wr, name, td)
}

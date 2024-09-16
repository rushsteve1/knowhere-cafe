package models

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"time"

	"knowhere.cafe/src/shared"
)

const (
	layout_key   = "templates/layout.html"
	fragment_key = "templates/fragment.html"
)

type TemplateState struct {
	tmpls map[string]*template.Template
}

func SetupTemplates(fs fs.FS, patterns ...string) TemplateState {
	slog.Info("compiling templates")

	state := TemplateState{make(map[string]*template.Template, 2+len(patterns))}

	funcMap := template.FuncMap{}

	patterns = append(patterns, layout_key, fragment_key)

	for _, pattern := range patterns {
		tmpl := template.New(pattern)
		tmpl = tmpl.Funcs(funcMap)
		state.tmpls[pattern] = template.Must(tmpl.ParseFS(fs, pattern))
	}

	return state
}

type TemplateData[T any] struct {
	Cfg      Config
	Now      time.Time
	PageName string
	Data     T
}

func RenderFullTemplate[T any](
	ctx context.Context,
	wr io.Writer,
	name string,
	data T,
) error {
	state := shared.Must(CtxState(ctx))
	cfg, err := state.Config()
	if err != nil {
		return err
	}

	td := TemplateData[T]{
		cfg, time.Now(), "", data,
	}

	return state.Templ.tmpls[layout_key].ExecuteTemplate(wr, name, td)
}

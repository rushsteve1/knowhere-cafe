package models

import (
	"context"
	"io"
	"strings"
	"time"

	"knowhere.cafe/src/shared/log"
)

type TemplateData[T any] struct {
	Cfg      Config
	Now      time.Time
	PageName string
	Data     T
}

func Render[T any](
	ctx context.Context,
	wr io.Writer,
	name string,
	data T,
) error {
	state := log.Must(CtxState(ctx))
	cfg, err := state.Config()
	if err != nil {
		return err
	}

	pageName := strings.TrimSuffix(name, ".html")

	td := TemplateData[T]{
		cfg, time.Now(), pageName, data,
	}

	return state.Templ.ExecuteTemplate(wr, name, td)
}

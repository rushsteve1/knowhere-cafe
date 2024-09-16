package web

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"os"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
)

//go:embed static
var staticFiles embed.FS

func StaticFiles(ctx context.Context) fs.FS {
	state := shared.Must(models.CtxState(ctx))
	if state.Flags.Dev {
		slog.WarnContext(ctx, "loading static files from fs")
		return os.DirFS("web/static")
	}
	return staticFiles
}

//go:embed templates
var templateFiles embed.FS

func TemplateFiles(ctx context.Context) fs.FS {
	state := shared.Must(models.CtxState(ctx))
	if state.Flags.Dev {
		// slog.WarnContext(ctx, "loading templates from fs")
		// return os.DirFS("templates")
	}
	return templateFiles
}

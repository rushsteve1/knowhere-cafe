// Embeds folders into the binary, and sets up functions for
// reading directly from those folders in Dev mode

package main

import (
	"embed"
	"io/fs"
	"os"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared/easy"
)

func StaticFiles(flags models.FlagConfig) (out fs.FS) {
	return devEmbed(flags, "static", staticFiles)
}

func TemplateFiles(flags models.FlagConfig) fs.FS {
	return devEmbed(flags, "templates", templateFiles)
}

func devEmbed(flags models.FlagConfig, path string, or embed.FS) (out fs.FS) {
	if flags.Dev {
		return os.DirFS(path)
	}
	return easy.Must(fs.Sub(or, path))
}

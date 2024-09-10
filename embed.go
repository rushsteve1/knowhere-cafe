package main

import "embed"

//go:embed static
var StaticFiles embed.FS

//go:embed templates
var TemplateFiles embed.FS

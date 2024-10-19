package models

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/lib/pq"
)

type ServerStartup struct {
	Time      time.Time       `gorm:"primaryKey"`
	BuildInfo debug.BuildInfo `gorm:"serializer:json"`
	Environ   pq.StringArray  `gorm:"type:text[]; not null; default:'{}'"`
	Cwd       string          `gorm:"not null"`
	ConfigID  uint            `gorm:"not null"`
	Config    Config
}

func NewServerStartup(cfg Config) ServerStartup {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("missing buildinfo")
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic("missing cwd")
	}

	return ServerStartup{
		Time:      time.Now(),
		ConfigID:  cfg.ID,
		BuildInfo: *buildInfo,
		Environ:   os.Environ(),
		Cwd:       cwd,
	}
}

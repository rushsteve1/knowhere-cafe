package models

import (
	"runtime/debug"
	"time"
)

type ServerStartup struct {
	Time      time.Time `gorm:"primaryKey"`
	ConfigID  uint      `gorm:"not null"`
	Config    Config
	BuildInfo debug.BuildInfo `gorm:"serializer:json"`
}

func NewServerStartup(cfg Config) ServerStartup {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("missing buildinfo")
	}

	return ServerStartup{
		Time:      time.Now(),
		ConfigID:  cfg.ID,
		BuildInfo: *buildInfo,
	}
}

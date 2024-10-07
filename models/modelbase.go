package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelBase struct {
	ID        uuid.UUID `gorm:"primarykey; type:uuid; default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m ModelBase) ETag() string {
	return formatTag(m.UpdatedAt)
}

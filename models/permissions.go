package models

import (
	"database/sql"
	"reflect"
	"time"

	"knowhere.cafe/src/shared"
)

type Permissions struct {
	ModelBase
	// This can be a role, group, or user
	AppliesTo     string `gorm:"uniqueIndex"`
	Administrator sql.NullBool
	CanPost       sql.NullBool
	DailyLimit    sql.NullInt16
	PostSpeed     sql.Null[time.Duration]
	CanHideEmail  sql.NullBool
	CanVote       sql.NullBool
	CanReport     sql.NullBool
}

func defaultAnonPermissions() Permissions {
	return Permissions{
		AppliesTo: shared.ANON_USER_ROLE,
	}
}

func defaultRegisteredPermissions() Permissions {
	return Permissions{
		AppliesTo:    shared.REGISTERED_USER_ROLE,
		CanPost:      sql.NullBool{Bool: true, Valid: true},
		DailyLimit:   sql.NullInt16{Int16: 3, Valid: true},
		PostSpeed:    sql.Null[time.Duration]{V: time.Hour, Valid: true},
		CanHideEmail: sql.NullBool{Bool: true, Valid: true},
		CanVote:      sql.NullBool{Bool: true, Valid: true},
		CanReport:    sql.NullBool{Bool: true, Valid: true},
	}
}

func FlattenPermissions(roles []string, perms []Permissions) (out Permissions) {
	// Sort the permissions into the same order as the roles
	for _, r := range roles {
		for _, p := range perms {
			if r == p.AppliesTo {
				// A bit of reflection or automatically grab all the fields
				vo := reflect.ValueOf(out)
				vp := reflect.ValueOf(p)
				for i, f := range reflect.VisibleFields(vo.Type()) {
					if f.Name == "AppliesTo" {
						continue
					}
					vo.Field(i).Set(vp.Field(i))
				}
			}
		}
	}
	return out
}

// Group of users
type Group struct {
	ModelBase
	Name   string `gorm:"index"`
	Bio    string `gorm:"default:''; not null"`
	Users  []User `gorm:"many2many:user_groups"`
	CanHat bool
}

package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex"`
	LastSeenAt time.Time
	Roles      []string
	InvitedBy  uint `gorm:"references:users(id)"`
}

type Permissions struct {
	AppliesTo    string `gorm:"primaryKey"`
	CanPost      bool
	PostLimit    int
	CanHideEmail bool
	CanVote      bool
	CanFlag      bool
}

type Invite struct {
	gorm.Model
	UserID    uint `gorm:"references:users(id)"`
	Limit     int
	ExpiresAt sql.NullTime
}

func (u User) AllRoles() []string {
	return append([]string{"registered", "anon", u.Email}, u.Roles...)
}

// WebAuthn
// https://pkg.go.dev/github.com/go-webauthn/webauthn/webauthn#User

func (u User) WebAuthnID() []byte {
	// TODO ensure this is valid and secure
	hasher := sha1.New()
	binary.Write(hasher, binary.NativeEndian, u.ID)
	return hasher.Sum(nil)
}

func (u User) WebAuthnName() string {
	return u.Email
}

func (u User) WebAuthnDisplayName() string {
	return u.WebAuthnName()
}

func (u User) WebAuthnCredentials() []webauthn.Credential {

}

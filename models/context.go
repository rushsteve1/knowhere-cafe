package models

import (
	"context"
	"net/smtp"
	"text/template"

	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"
	"knowhere.cafe/src/shared"
)

type ContextState struct {
	Flags FlagConfig
	DB    *gorm.DB
	IMAP  *imapclient.Client
	SMTP  *smtp.Client
	Templ *template.Template
}

func (cs ContextState) Config() (Config, error) {
	var cfg Config
	res := cs.DB.Last(&cfg)
	return cfg, res.Error
}

func CtxState(ctx context.Context) (ContextState, error) {
	cs, ok := ctx.Value(shared.CTX_STATE_KEY).(ContextState)
	if !ok {
		return ContextState{}, shared.ErrMissingState{}
	}
	return cs, nil
}

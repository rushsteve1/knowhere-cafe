package models

import (
	"context"
	"net/smtp"

	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"
	"knowhere.cafe/src/shared"
)

type ContextState struct {
	Flags FlagConfig
	DB    *gorm.DB
	IMAP  *imapclient.Client
	SMTP  *smtp.Client
	Templ TemplateState
}

func (cs ContextState) Config() (Config, error) {
	var cfg Config
	res := cs.DB.Last(&cfg)
	return cfg, res.Error
}

func State(ctx context.Context) (ContextState, error) {
	cs, ok := ctx.Value(shared.CTX_STATE_KEY).(ContextState)
	if !ok {
		return ContextState{}, shared.ErrMissingState
	}
	return cs, nil
}

package models

import (
	"context"

	"gorm.io/gorm"
	"tailscale.com/client/tailscale/apitype"
	"tailscale.com/tsnet"

	"knowhere.cafe/src/shared"
)

type ContextState struct {
	Flags FlagConfig
	DB    *gorm.DB
	Templ *TemplateState
	Tsnet *tsnet.Server
}

func (cs ContextState) Config() (Config, error) {
	var cfg Config
	res := cs.DB.Last(&cfg)
	return cfg, res.Error
}

const STATE_CTX_KEY = "STATE"

// TODO change this singature to an ok bool instead of err
func State(ctx context.Context) (state ContextState, err error) {
	cs, ok := ctx.Value(STATE_CTX_KEY).(ContextState)
	if !ok {
		return ContextState{}, shared.ErrMissingState
	}
	return cs, nil
}

const AUTH_CTX_KEY = "AUTH"

type ContextAuth *apitype.WhoIsResponse

func Auth(ctx context.Context) (who *ContextAuth, err error) {
	who, ok := ctx.Value(AUTH_CTX_KEY).(*ContextAuth)
	if !ok {
		return nil, shared.ErrNotAuth
	}
	return who, nil
}

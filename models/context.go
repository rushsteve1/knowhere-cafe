package models

import (
	"context"

	"gorm.io/gorm"
	"knowhere.cafe/src/shared"
)

type ContextState struct {
	Flags FlagConfig
	DB    *gorm.DB
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

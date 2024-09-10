// Tools for working with [context.Context]

package shared

import (
	"context"
	"net/smtp"

	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"
	"knowhere.cafe/src/models"
)

type ContextState struct {
	FlagCfg *models.FlagConfig
	Cfg     *models.Config
	DB      *gorm.DB
	IMAP    *imapclient.Client
	SMTP    *smtp.Client
}

func CtxData(ctx context.Context) ContextState {
	return ctx.Value(CTX_STATE_KEY).(ContextState)
}

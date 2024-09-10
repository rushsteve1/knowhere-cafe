// Tools for working with [context.Context]

package shared

import (
	"context"
	"net/smtp"

	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"
	"knowhere.cafe/src/models"
)

type ContextData struct {
	FlagCfg *models.FlagConfig
	Cfg     *models.Config
	DB      *gorm.DB
	IMAP    *imapclient.Client
	SMTP    *smtp.Client
}

func CtxData(ctx context.Context) ContextData {
	return ctx.Value(CTX_DATA_KEY).(ContextData)
}

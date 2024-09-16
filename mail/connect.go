package mail

import (
	"net/smtp"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"

	"github.com/emersion/go-imap/v2/imapclient"
)

func IMAPConnect(cred models.ConfigCredentials) (*imapclient.Client, error) {
	return nil, shared.ErrUnimplemented{}
}

func SMTPConnect(cred models.ConfigCredentials) (*smtp.Client, error) {
	return nil, shared.ErrUnimplemented{}
}

package client

import (
	"github.com/lancer-kit/sender/models/email"
	"github.com/lancer-kit/sender/models/sms"
)

type Client interface {
	SendSms(sms.Message) error
	SendEmail(email.Message) error
}

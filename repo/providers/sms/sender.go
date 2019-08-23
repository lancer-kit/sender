package sms

import (
	"github.com/lancer-kit/sender/models/sms"
)

type Sender interface {
	SendSms(sms *sms.Message) error
}

package email

import (
	"github.com/lancer-kit/sender/models/email"
)

type Sender interface {
	SendEmail(email *email.Message) error
}

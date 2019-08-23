package email

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Topic is a topic in async channel, through which the sender receives new messages.
const Topic = "sender.emails"

type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.To, validation.Required, is.Email),
		validation.Field(&m.Subject, validation.Required),
		validation.Field(&m.Body, validation.Required),
	)
}

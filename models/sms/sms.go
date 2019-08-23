package sms

import validation "github.com/go-ozzo/ozzo-validation"

// Topic is a topic in async channel, through which the sender receives new otp messages.
const Topic = "sender.sms"

type Provider string

const (
	ProviderViber    Provider = "viber"
	ProviderWhatsApp Provider = "whatsapp"
	ProviderSMS      Provider = "sms"
	ProviderTelegram Provider = "telegram"
)

type Message struct {
	// Provider indicates which service to use to send SMS.
	Provider Provider `json:"provider,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Text     string   `json:"text"`
}

func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Phone, validation.Required),
		validation.Field(&m.Text, validation.Required),
		validation.Field(&m.Provider, validation.Required, validation.In(
			ProviderSMS,
			ProviderTelegram,
			ProviderViber,
			ProviderWhatsApp,
		)),
	)
}

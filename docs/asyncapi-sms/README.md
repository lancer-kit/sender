# AsyncAPI-SMS

`asyncapi-sms` is a NATS API that is used to send SMSs.

## Usage

You need to use package `sender_client` in order to communicate with sender.
To send an SMS using `asyncapi-sms`, just send `SmsInput` structure (marshaled into JSON) to `sender_client.SmsTopic` NATS channel.

#### type SmsInput

```go
type SmsInput struct {
	// TTL is number of seconds during which object will live in cache and be throttled. Zero means no throttling
	TTL int `json:"ttl"`
	// Phone indicates the phone number to send SMS to.
	Phone string `json:"phone,omitempty"`
	// Type indicates which template will be used.
	TemplateId string `json:"template_id"`
	// Provider indicates which service to use to send SMS.
	Provider Provider `json:"provider,omitempty"`
	// Data to fill in the template, depends on the `Type`.
	Data map[string]string `json:"data"`
	// If not empty - Data is ignored and body of sms is Raw string
	Raw string `json:"raw"`
}
```

#### type Provider

```go
type Provider int

const (
	ViberProvider Provider = 1 + iota
	WAProvider
	SMSProvider
	TelegramProvider
)
```

`Provider` – enum of SMS providers.

`ViberProvider` – Viber.

`WAProvider` – WhatsApp.

`SMSProvider` – SMS.

`TelegramProvider` – Telegram.
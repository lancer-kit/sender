# AsyncAPI-Email

`asyncapi-email` is a NATS API that is used to send emails.

## Usage

You need to use package `sender_client` in order to communicate with sender.
To send an email using `asyncapi-email`, just send `EmailInput` structure (marshaled into JSON) to `sender_client.EmailTopic` NATS channel.

#### type EmailInput

```go
type EmailInput struct {
	// TTL is number of seconds during which object will live in cache and be throttled. Zero means no throttling
	TTL int `json:"ttl"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	// Type indicates which template will be used.
	TemplateId string `json:"template_id"`
	// Data to fill in the template, depends on the `Type`.
	Data map[string]string `json:"data"`
	// If not empty - Data is ignored and body of email is Raw string
	Raw string `json:"raw"`
}
```
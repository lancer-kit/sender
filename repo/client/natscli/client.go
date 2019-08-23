package natscli

import (
	"encoding/json"

	"github.com/lancer-kit/sender/models/email"
	"github.com/lancer-kit/sender/models/sms"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
)

type client struct {
	conn *nats.Conn
	url  string
}

func New(url string) (c *client, err error) {
	c = &client{
		url: url,
	}

	if err = c.ensure(); err != nil {
		return nil, errors.Wrap(err, "failed nats ensure")
	}

	return c, nil
}

func (c client) SendEmail(msg email.Message) (err error) {
	if err = msg.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	if err = c.publishJSON(email.Topic, msg); err != nil {
		return errors.Wrap(err, "publish failed")
	}

	return nil
}

func (c client) SendSms(msg sms.Message) (err error) {
	if err = msg.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	if err = c.publishJSON(sms.Topic, msg); err != nil {
		return errors.Wrap(err, "publish failed")
	}

	return nil
}

func (c *client) ensure() error {
	var err error

	if c.url == "" {
		return errors.New("nats: connection url is not set")
	}

	if c.conn, err = nats.Connect(c.url); err != nil {
		return errors.Wrap(err, "nats: failed to connect")
	}

	return nil
}

func (c *client) publishJSON(topic string, msg interface{}) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.conn.Publish(topic, raw)
}

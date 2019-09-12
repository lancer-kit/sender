package restcli

import (
	"fmt"
	"net/http"

	"github.com/lancer-kit/armory/api/httpx"
	"github.com/lancer-kit/sender/models/email"
	"github.com/lancer-kit/sender/models/sms"
	"github.com/lancer-kit/sender/repo/client"
	"github.com/pkg/errors"
)

type cli struct {
	conn httpx.Client
	url  string
}

func New(url string) (c client.Client) {
	return &cli{
		url:  url,
		conn: httpx.NewXClient(),
	}
}

func (c cli) SendEmail(msg email.Message) (err error) {
	if err = msg.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	var res *http.Response
	res, err = c.conn.PostJSON(fmt.Sprintf("%s/v1/sender/email", c.url), msg, nil)
	if err != nil {
		return errors.Wrap(err, "publish failed")
	}
	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	return nil
}

func (c cli) SendSms(msg sms.Message) (err error) {
	if err = msg.Validate(); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	var res *http.Response
	res, err = c.conn.PostJSON(fmt.Sprintf("%s/v1/sender/sms", c.url), msg, nil)
	if err != nil {
		return errors.Wrap(err, "publish failed")
	}
	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	return nil
}

package viber

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/sms"
	smsp "github.com/lancer-kit/sender/repo/providers/sms"
	"github.com/pkg/errors"
)

type sender struct {
	cfg    *config.Viber
	client *http.Client
}

func New(cfg *config.Viber) smsp.Sender {
	return &sender{
		cfg:    cfg,
		client: new(http.Client),
	}
}

func (s *sender) SendSms(msg *sms.Message) error {
	req := s.makeRequest(msg)

	res, err := s.client.Do(req)

	if err != nil {
		return errors.Wrap(err, "cant do request to Viber")
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("viber respond with status %s", res.Status)
	}

	return nil
}

func (s *sender) makeRequest(msg *sms.Message) *http.Request {
	body := url.Values{}
	body.Add("body", msg.Text)
	body.Add("phone", msg.Phone)
	body.Add("messenger", "viber")
	body.Add("project_id", "0")
	rb := *strings.NewReader(body.Encode())

	req, _ := http.NewRequest(http.MethodPost, s.cfg.APIURL, &rb)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

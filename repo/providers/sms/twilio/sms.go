package twilio

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
	cfg    *config.Twilio
	client *http.Client
}

func New(cfg *config.Twilio) smsp.Sender {
	return &sender{
		cfg:    cfg,
		client: new(http.Client),
	}
}

func (s *sender) SendSms(msg *sms.Message) error {
	req := s.makeRequest(msg)
	res, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant do request to Twilio")
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("twillio respond with status %s", res.Status)
	}

	return nil
}

func (s *sender) makeRequest(msg *sms.Message) *http.Request {
	body := url.Values{}
	body.Set("To", msg.Phone)
	body.Set("From", s.cfg.Sender)
	body.Set("Body", msg.Text)
	rb := *strings.NewReader(body.Encode())

	req, _ := http.NewRequest(http.MethodPost, s.cfg.APIURL, &rb)
	req.SetBasicAuth(s.cfg.AccountSid, s.cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req
}

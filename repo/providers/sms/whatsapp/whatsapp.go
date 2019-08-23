package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lancer-kit/sender/config"
	"github.com/lancer-kit/sender/models/sms"
	smsp "github.com/lancer-kit/sender/repo/providers/sms"
	"github.com/pkg/errors"
)

type sender struct {
	cfg *config.Whatsapp
}

func New(cfg *config.Whatsapp) smsp.Sender {
	return &sender{
		cfg: cfg,
	}
}

func (s *sender) SendSms(msg *sms.Message) error {
	body := make(map[string]string)
	body["phone"] = msg.Phone
	body["body"] = msg.Text

	b, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "cannot marshal request body")
	}

	res, err := http.Post(s.cfg.APIURL+s.cfg.APIKey, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "cant do request to WhatsApp")
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("whatsapp respond with status %s", res.Status)
	}

	return nil
}

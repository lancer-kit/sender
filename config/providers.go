package config

import "github.com/pkg/errors"

type Providers struct {
	Email Email `yaml:"email"`
	SMS   SMS   `yaml:"sms"`
}

func (p Providers) Validate() (err error) {
	if err = p.Email.Validate(); err != nil {
		return errors.Wrap(err, "email")
	}
	if err = p.SMS.Validate(); err != nil {
		return errors.Wrap(err, "sms")
	}

	return nil
}

type Email struct {
	Mailgun  *Mailgun  `yaml:"mailgun"`
	Sendgrid *Sendgrid `yaml:"sendgrid"`
	SMTP     *SMTP     `yaml:"smtp"`
}

func (e Email) Validate() (err error) {
	if e.Mailgun != nil {
		err = e.Mailgun.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "mailgun")
	}

	if e.Sendgrid != nil {
		err = e.Sendgrid.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "sendgrid")
	}

	if e.SMTP != nil {
		err = e.SMTP.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "smtp")
	}
	return nil
}

type SMS struct {
	Whatsapp *Whatsapp `yaml:"whatsapp"`
	Twilio   *Twilio   `yaml:"twilio"`
	Viber    *Viber    `yaml:"viber"`
}

func (s SMS) Validate() (err error) {
	if s.Whatsapp != nil {
		err = s.Whatsapp.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "mailgun")
	}

	if s.Twilio != nil {
		err = s.Twilio.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "sendgrid")
	}

	if s.Viber != nil {
		err = s.Viber.Validate()
	}
	if err != nil {
		return errors.Wrap(err, "smtp")
	}
	return nil
}

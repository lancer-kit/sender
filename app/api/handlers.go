package api

import (
	"net/http"

	"github.com/lancer-kit/armory/api/httpx"
	"github.com/lancer-kit/armory/api/render"
	"github.com/lancer-kit/armory/log"
	"github.com/lancer-kit/sender/models/email"
	"github.com/lancer-kit/sender/models/sms"
	emailp "github.com/lancer-kit/sender/repo/providers/email"
	smsp "github.com/lancer-kit/sender/repo/providers/sms"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger      *logrus.Entry
	emailSender emailp.Sender
	smsSenders  map[sms.Provider]smsp.Sender
}

func (h handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = new(email.Message)
	)

	if err = httpx.ParseJSONBody(r, &msg); err != nil {
		render.BadRequest(w, "invalid json")
		return
	}

	if err = msg.Validate(); err != nil {
		h.log(r).WithError(err).Debugln("validation failed")
		render.BadRequest(w, "json validation failed")
		return
	}

	if err = h.emailSender.SendEmail(msg); err != nil {
		h.log(r).WithError(err).Error("cannot send email")
		render.ServerError(w)
		return
	}

	render.Success(w, "message was sent")
}

func (h handler) SendSms(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		msg = new(sms.Message)
	)

	if err = httpx.ParseJSONBody(r, &msg); err != nil {
		render.BadRequest(w, "invalid json")
		return
	}

	if err = msg.Validate(); err != nil {
		render.BadRequest(w, "json validation failed")
		return
	}

	if _, ok := h.smsSenders[msg.Provider]; !ok {
		h.log(r).WithError(err).Error("cannot send sms")
		render.BadRequest(w, "unavailable provider")
		return
	}

	if err = h.smsSenders[msg.Provider].SendSms(msg); err != nil {
		h.log(r).WithError(err).Error("cannot send sms")
		render.ServerError(w)
		return
	}

	render.Success(w, "sms was sent")
}

func (h handler) log(r *http.Request) *logrus.Entry {
	return log.IncludeRequest(h.logger, r)
}

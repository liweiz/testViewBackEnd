package testView

import (
	"bytes"
	"errors"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/smtp"
	"path"
	"text/template"
)

// emailCategory
const (
	emailForActivation = iota
	emailForPasswordResetting
)

func GenerateEmail(emailCategory int, uniqueUrl string, recipient string) (err error) {
	var emailFileName string
	d := &injectedData{"matt.z.lw@gmail.com", recipient, uniqueUrl}
	switch emailCategory {
	case emailForActivation:
		emailFileName = "email_activation.html"
	case emailForPasswordResetting:
		emailFileName = "email_password_resetting.html"
	default:
		err = errors.New("No email file set for this purpose.")
		return
	}
	emailFilePath := path.Join("pages/emails/", emailFileName)
	var t *template.Template
	t, err = template.ParseFiles("pages/layout.html", emailFilePath)
	if err == nil {
		var b bytes.Buffer
		err = t.Execute(&b, d)
		if err == nil {
			auth := smtp.PlainAuth(
				"",
				"matt.z.lw@gmail.com",
				"2Q_b03rFhPSKqIbWjYaOVA",
				"smtp.mandrillapp.com:587")
			err = smtp.SendMail(recipient, auth, d.From, recipient, b)
		}
	}
	return
}

type injectedData struct {
	From      string
	To        string
	UniqueUrl string
}

package testView

import (
	"bytes"
	"errors"
	"github.com/go-martini/martini"
	// "github.com/twinj/uuid"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/smtp"
	"path"
	"text/template"
)

// emailCategory
const (
	EmailForActivation = iota
	EmailForPasswordResetting
)

/*
   No request body since it's a GET call from client.
   The whole activation process:
   1. User presses the activation url in the activation email sent to the user's email. Only need to generate url once since activation is a one time activity.
   2. A webpage is shown in user's browser. If it is not activated before, the html shows message: account activated. Otherwise, the web page shows message: account has been activated already.
   3. The clients will update the activated state through the next sync request.
*/

/*
   No request body since it's a GET call from client.

   The whole resetting process:
   1. User presses the password resetting button in the activation email sent to the user's email.
   2. A webpage is shown in user's browser. If the link is not valid, the html shows the interface to reset the password. Otherwise, the web page shows message: invalid link, please require a new link by pressing the resetting button in your app.
*/

func EmailSender(emailCategory int) martini.Handler {
	return func(db *mgo.Database, rw http.ResponseWriter, req *http.Request, logger *log.Logger, params martini.Params) {
		userId := bson.ObjectIdHex(params["user_id"])
		fmt.Println("0", userId.Hex())
		err := SendEmail(emailCategory, db, req, params, userId)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
		} else {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

func SendEmail(emailCategory int, db *mgo.Database, req *http.Request, params martini.Params, userId bson.ObjectId) (err error) {
	var aUser User
	err = db.C("users").Find(bson.M{"_id": userId}).One(&aUser)
	if err == nil {
		switch emailCategory {
		case EmailForActivation:
			// Send an email with the activation link. E.g., http://www.xxx.com/users/:user_id/activation/:activation_code
			uniqueUrl := path.Join(req.URL.String(), aUser.ActivationUrlCode)
			err = GenerateEmail(emailCategory, uniqueUrl, aUser.Email)
		case EmailForPasswordResetting:
			_, err = UpdateNonDicDB(UserAddPasswordUrlCode, nil, db, params, userId)
			if err == nil {
				// user has been updated, get the latest version.
				var bUser User
				err = db.C("users").Find(bson.M{"_id": userId}).One(&bUser)
				if err == nil {
					l := len(bUser.PasswordResettingUrlCodes) - 1
					if l >= 0 {
						uniqueUrlCode := bUser.PasswordResettingUrlCodes[l].PasswordResettingUrlCode
						// "/users/:user_id/password/:password_resetting_code"
						uniqueUrl := path.Join("/users", userId.Hex(), "password", uniqueUrlCode)
						err = GenerateEmail(emailCategory, uniqueUrl, aUser.Email)
						if err == nil && l > 0 {
							_, _ = UpdateNonDicDB(UserRemovePasswordUrlCodeRoutinely, nil, db, params, userId)
						}
					} else {
						err = errors.New("Empty PasswordResettingUrlCodes array, previous insertion failed.")
					}
				} else {
					fmt.Println("2", err.Error())
				}
			}
		}
	} else {
		fmt.Println("1", err.Error())
	}
	return
}

func GenerateEmail(emailCategory int, uniqueUrl string, recipient string) (err error) {
	var emailFileName string
	var s string
	d := &injectedData{"http://localhost:3000", "matt.z.lw@gmail.com", recipient, uniqueUrl}
	switch emailCategory {
	case EmailForActivation:
		emailFileName = "email_activation.html"
		s = "Subject: Activate your Remolet account and you will be all set.\r\n"
	case EmailForPasswordResetting:
		emailFileName = "email_password_resetting.html"
		s = "Subject: Reset your Remolet password\r\n"
	default:
		err = errors.New("No email file set for this purpose.")
		return
	}
	emailFilePath := path.Join("pages/emails/", emailFileName)
	var t *template.Template
	t, err = template.ParseFiles("pages/layout/email_layout.html", emailFilePath)
	if err == nil {
		var b bytes.Buffer
		err = t.Execute(&b, d)
		if err == nil {
			auth := smtp.PlainAuth(
				"",
				"remolet.z@gmail.com",
				"byeoWdNChHJZEQejH7pbpg",
				"smtp.mandrillapp.com")
			f := "From: remolet.z@gmail.com\r\n"
			m := "MIME-Version: 1.0\r\n"
			c := "Content-Type: text/html\r\n"
			// e := "Content-Transfer-Encoding: Base64\r\n"
			msg := f + s + m + c + string(b.Bytes())
			err = smtp.SendMail("smtp.mandrillapp.com:587", auth, d.From, []string{recipient}, []byte(msg))
		}
	}
	return
}

type injectedData struct {
	MyHost    string
	From      string
	To        string
	UniqueUrl string
}

package testView

import (
	"errors"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"time"
)

func UrlCodeChecker() martini.Handler {
	return func(db *mgo.Database, params martini.Params, logger *log.Logger, rw http.ResponseWriter) {
		err := CheckUrlCode(db, params)
		if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

// Check if the url is valid. Each url is only valid for one hour and one time use.
func CheckUrlCode(db *mgo.Database, params martini.Params) (err error) {
	userId := bson.ObjectIdHex(params["user_id"])
	urlCode := params["password_resetting_code"]

	var aUser User
	err = db.C("users").Find(bson.M{"_id": userId}).One(&aUser)
	if err == nil {
		if len(aUser.PasswordResettingUrlCodes) > 0 {
			for _, u := range aUser.PasswordResettingUrlCodes {
				if u.PasswordResettingUrlCode == urlCode {
					if time.Now().UnixNano()-u.TimeStamp < (3.6e+12) {
						return
					}
				}
			}
		}
	}
	if err == nil {
		err = errors.New("Previous change was successful or the link is expired.")
	}
	return
}

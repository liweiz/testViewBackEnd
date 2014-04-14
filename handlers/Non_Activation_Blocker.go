package testView

import (
	"errors"
	// "fmt"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

// Only requests that pass the gateKeeper are processed here. This indicates a not found err here means user not activated.
func NonActivationBlocker() martini.Handler {
	return func(db *mgo.Database, params martini.Params, logger *log.Logger, rw http.ResponseWriter) {
		err := BlockNonActivationReq(db, params)
		if err == mgo.ErrNotFound {
			err = errors.New("User not activated.")
		}
		if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

func BlockNonActivationReq(db *mgo.Database, params martini.Params) (err error) {
	userId := bson.ObjectIdHex(params["user_id"])
	var aUser User
	err = db.C("users").Find(bson.M{"_id": userId, "activated": true}).One(&aUser)
	return
}

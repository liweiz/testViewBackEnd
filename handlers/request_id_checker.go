package testView

import (
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

// Two modules needed for reqVerNo check.
// 1. One module in the beginning. Check if it has been processed. If yes, return a 200 response. If not, go ahead.
// 2. The other in the end. If it is successfully processed, add to the list. Only successful ones in the list.

func ReqIdChecker() martini.Handler {
	return func(db *mgo.Database, v *Vehicle, params martini.Params, rw martini.ResponseWriter, logger *log.Logger) {
		isProcessed, err := CheckReqId(db, v.RequestId, v.DeviceUUID, bson.ObjectIdHex(params["user_id"]))
		if !isProcessed {
			if err == mgo.ErrNotFound {
				// Good to proceed
			} else if err != nil {
				WriteLog(err.Error(), logger)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		} else {
			rw.WriteHeader(http.StatusOK)
		}
	}
}

func CheckReqId(db *mgo.Database, reqIdToCheck string, deviceUUID string, userId bson.ObjectId) (isProcessed bool, err error) {
	_, err = db.C("requestProcessed").Find(bson.M{
		"belongToUser": userId,
		"requestId":    reqIdToCheck,
		"deviceUUID":   deviceUUID}).Count()
	if err == nil {
		isProcessed = true
	}
	return
}

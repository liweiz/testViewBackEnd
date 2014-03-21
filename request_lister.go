package testView

import (
		"net/http"
		"encode/json"
		"github.com/codegangsta/martini"
		"labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "net/url"
        "error"
        "reflect"
)

// Two modules needed for reqVerNo check.
// 1. One module in the beginning. Check if it has been processed. If yes, return a 200 response. If not, go ahead.
// 2. The other in the end. If it is successfully processed, add to the list. Only successful ones in the list. 

func CheckIfAlreadyProcessed(params *martini.Params, ctx martini.Context, db *mgo.Database, rw martini.ResponseWriter) {
	clientUUID := GetVehicleContentInContext(ctx, "ClientUUID").String()
	reqVerNo := (int)(GetVehicleContentInContext(ctx, "ReqVerNo").Int())
	isProcessed, r := CheckReqProcessed(params, clientUUID, reqVerNo, db)
	if isProcessed {
		// Send 200 response
		rw.WriteHeader(StatusOK)
	}
	// Hand to next module, do nothing here.
}

func MarkSuccessfulProcessToReqList(params *martini.Params, ctx martini.Context, db *mgo.Database) {
	clientUUID := GetVehicleContentInContext(ctx, "ClientUUID").String()
	reqVerNo := (int)(GetVehicleContentInContext(ctx, "ReqVerNo").Int())
	AddToReqList(params, db, clientUUID, reqVerNo)
	// Hand to next module, do nothing here.
}



func CheckReqProcessed(params *martini.Params, clientUUID string, reqVerNo int, db *mgo.Database) (isProcessed bool) {
	err := db.C("RequestProcessed").Find(bson.M{
		"deviceUUID": clientUUID,
		"belongToUser": params["user_id"]，
		"requestVersionNo": reqVerNo
		}).One(r)
	if err != nil {
		isProcessed = false
	} else {
		isProcessed = true
	}
	return
}

func AddToReqList(params *martini.Params, db *mgo.Database, clientUUID string, reqVerNo int) {
	err := db.C("RequestProcessed").Insert(bson.M{
		"_id": bson.NewObjectId(),
		"belongToUser": bson.ObjectIdHex(params["user_id"]),
		"deviceUUID": clientUUID,
		"requestVersionNo": reqVerNo,
		"timestamp": time.Now().UnixNano()
		})
	if err != nil {
		// 
	}
}
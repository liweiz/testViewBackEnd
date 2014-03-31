package testView

import (
	"errors"
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
)

const (
	// For new card only
	ConflictCardAlreadyExists = iota
	// conflict detail and card(s), possiblely to have two cards, on server that causes the conflict
	ConflictOverwriteClient
	NoConflictOverwriteDB
	ConflictCreateAnotherInDB
)

func HandleCardVersionConflict(action string, docFromReq interface{}, docFromDB interface{}) (decisionCode int) {
	// Assume that the uniqueness check has been done before this. So every docFromReq is unique according to the db.
	d := reflect.ValueOf(docFromDB)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	r := reflect.ValueOf(docFromReq)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	dv := d.FieldByName("versionNo").Int()
	rv := r.FieldByName("versionNo").Int()
	// Currently for card only
	// docName := reflect.TypeOf(docFromDB).Name()
	if dv > rv {
		// A newer version is on server already. For updating, create the one from client as a new doc if it's unique. In the case that the card on server has been deleted, which also indicates the card on server has a higher versionNo, create the updating card from client as a new card instead of changing the isDeleted field. Because this is consistent with the fundamental principle above:
		// "A newer version is on server already. For updating, create the one from client as a new doc if it's unique.
		// VersionNo is the only criterion.
		// For deleting, ignore the delete action and update the doc on client.
		// Creating is not operate on existing docs, so no need to compare versionNo, only need to check uniqueness.
		if action == "update" {
			decisionCode = ConflictCreateAnotherInDB
		} else if action == "delete" {
			// If card on server has been deleted, the delete action has no effect. The server will ask the client to delete the card on server as well. The delete request sent by client just happens to be the same decision made  by the server and the version comparison leads to overwriting the client.
			// If the card on server is not deleted, the higher version f the card on server indicates the client needs to comply with the server's dicision.
			// In both cases, we need to overwrite card on client.
			decisionCode = ConflictOverwriteClient
		}
	} else if dv == rv {
		decisionCode = NoConflictOverwriteDB
	} else if dv < rv {
		// Impossible. If this does happens, overwrite the client anyway.
		decisionCode = ConflictOverwriteClient
	}
	return
}

func SetDeviceTokens(db *mgo.Database, structFromReq interface{}, params martini.Params) (tokens TokensInCommon, err error) {
	x := reflect.ValueOf(structFromReq)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	var newId bson.ObjectId
	if params["user_id"] == "" {
		// SignUp, must create a new deviceTokens in db.
		newId, err = InsertNonDicDB(DeviceTokensNew, structFromReq, db, params)
		if err == nil {
			err = db.C("deviceTokens").Find(bson.M{"_id": newId}).Select(SelectDeviceTokensInCommon).One(&tokens)
		}
	} else {
		// Get it from url.
		userId := bson.ObjectIdHex(params["user_id"])
		if err == nil {
			selector := bson.M{
				// UUID is available in both ReqSignUpOrIn and ReqRenewTokens.
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"belongTo":   userId}
			_, err = db.C("deviceTokens").Find(selector).Count()
			if err == nil {
				var selector1 bson.M
				selector1, err = UpdateNonDicDB(DeviceTokensUpdateTokens, structFromReq, db, params)
				if err == nil {
					err = db.C("deviceTokens").Find(selector1).Select(SelectDeviceTokensInCommon).One(&tokens)
				}
			} else {
				newId, err = InsertNonDicDB(DeviceTokensNew, structFromReq, db, params)
				if err == nil {
					err = db.C("deviceTokens").Find(bson.M{"_id": newId}).Select(SelectDeviceTokensInCommon).One(&tokens)
				}
			}
		}
	}
	return
}

// If card is unique, return true
func CheckCardUniqueness(db *mgo.Database, params martini.Params, cardValueInReq reflect.Value) (isUnique bool, duplicated []CardInCommon, err error) {
	x := reflect.ValueOf(cardValueInReq)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	idToCheck := bson.ObjectIdHex(params["user_id"])
	// Return the possible cards that duplicated for further process.
	err = db.C("cards").Find(bson.M{
		"belongTo": idToCheck,
		// Only check those not yet deleted. So it is possible to have multiple deleted cards with exactly the same content.
		"isDeleted":   false,
		"target":      x.FieldByName("Target").String(),
		"translation": x.FieldByName("Translation").String(),
		"context":     x.FieldByName("Context").String(),
		"detail":      x.FieldByName("Detail").String()}).All(&duplicated)
	if err == mgo.ErrNotFound {
		isUnique = true
	}
	return
}

// NEED TO REWRITE DUE TO DATA SCHEME CHANGE
// If cardInDic is unique, return true
// func CheckCardInDicUniqueness(db *mgo.Database, cardToCheck *CardInCommon) (isUnique bool, err error) {
// 	var result CardInDic
// 	err = db.C("cardInDics").Find(bson.M{
// 		// Only check those not yet deleted. So it is possible to have multiple deleted cards with exactly the same content.
// 		"isDeleted":   false,
// 		"isHidden":    false,
// 		"target":      cardToCheck.Target,
// 		"translation": cardToCheck.Translation,
// 		"context":     cardToCheck.Context,
// 		"detail":      cardToCheck.Detail}).One(&result)
// 	if err == mgo.ErrNotFound {
// 		isUnique = true
// 	}
// 	return
// }

// Default setting is the latest setting user uses. So compare all the deveiceInfo the user has and get the latest updated one to send to the new device as the initial settings on the device. Store it in db as well.
func GetDefaultDeviceInfo(userId bson.ObjectId, db *mgo.Database) (info *DeviceInfoInCommon, err error) {
	var r []DeviceInfo
	err = db.C("deviceInfos").Find(bson.M{
		"belongTo": userId}).All(&r)
	if err == nil {
		info, err = GetLatestUpdatedDeviceInfo(r)
	}
	return
}

func GetLatestUpdatedDeviceInfo(l []DeviceInfo) (info *DeviceInfoInCommon, err error) {
	if len(l) > 0 {
		r := FormIntSliceFromDocSlice(l, "LastModified")
		var x int64
		x, err = MaxInt(r)
		if err == nil {
			for i := range l {
				if l[i].LastModified == x {
					info.Id = l[i].Id
					info.BelongTo = l[i].BelongTo
					info.DeviceUUID = l[i].DeviceUUID
					info.SourceLang = l[i].SourceLang
					info.TargetLang = l[i].TargetLang
					info.SortOption = l[i].SortOption
					info.IsLoggedIn = l[i].IsLoggedIn
					info.RememberMe = l[i].RememberMe
					return
				}
			}
		}
	} else {
		err = errors.New("No element found in slice.")
	}
	return
}

func FormIntSliceFromDocSlice(docSlice interface{}, fieldName string) (r []int64) {
	d := reflect.ValueOf(docSlice)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	r = make([]int64, 0)
	if d.Len() > 0 {
		for i := 0; i < d.Len(); i++ {
			dd := d.Index(i)
			if dd.Kind() == reflect.Ptr {
				dd = dd.Elem()
			}
			r = append(r, dd.FieldByName(fieldName).Int())
		}
	}
	return
}

func MaxInt(ints []int64) (x int64, err error) {
	if len(ints) == 0 {
		err = errors.New("No Int found in slice.")
	} else {
		x = ints[0]
		if len(ints) > 1 {
			for i := range ints {
				if i < len(ints)-1 {
					if x < ints[i+1] {
						x = ints[i+1]
					}
				}
			}
		}
	}
	return
}

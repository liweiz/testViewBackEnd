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

func ProcessRequest(db *mgo.Database, route int, criteria bson.M, structFromReq interface{}, req *http.Request, structForRes interface{}, params *martini.Params) (err error) {
	m := req.Method
	// x and v are struct corresponding to JSON, so, to get the parts we need, there's one step further needed.
	x := reflect.ValueOf(structFromReq)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	v := reflect.ValueOf(structForRes)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if m == "GET" {
		// No potential conflict
		if route == oneCard {
			var resultCard CardInCommon
			err = db.C("cards").Find(criteria).Select(SelectCardInCommon).One(&resultCard)
			if !err {
				err = SetResBodyPart(v, "Cards", []CardInCommon{resultCard})
			}
		} else if route == oneUser {
			var resultUser UserInCommon
			err = db.C("users").Find(criteria).Select(SelectUserInCommon).One(&resultUser)
			if !err {
				err = SetResBodyPart(v, "User", resultUser)
			}
		} else if {
			var resultDeviceInfo DeviceInfoInCommon
			err := db.C("deviceInfos").Find(bson.M{
				"belongTo": params["user_id"],
				// /users/:user_id/devices/:device_id
				"_id": params["device_id"]}).Select(SelectDeviceInfoInCommon).One(&resultDeviceInfo)
			if !err {
				err = SetResBodyPart(v, "DeviceInfo", resultDeviceInfo)
			}
		}
	} else if m == "POST" {
		switch route {
		// Sign up
		/*
		SignUp flow:
		1. Client sends email and password to server
		2. Create user in db if everything's ok. Otherwise, return error message
		3. Client receives the user, lets the user set languagePair and sends it as deviceInfo to server.
		4. Server creates deviceInfo in db if everything's ok. Otherwise, return error message.

		SignIn flow:
		1. Client sends email and password to server
		2. Get user in db if everything's ok. Otherwise, return error message.
		2.1 No previous data on client found. Start sync process. If no deviceInfo on client, use the default one on server.
		2.2 Previous data on client. Start sync process.

		So a new user's deviceInfo is created on client by user and stored to db. A existing user's deviceInfo is delivered by sync process or providing the list for client to choose. When there is only one existing deviceInfo on server, send this to client.
		*/
		case signUp:
			// var aUser User
			// Check if already in use
			var n int
			n, err = db.C("users").Find(criteria).Count()
			if n > 0 && !err {
				var newId bson.ObjectId
				newId, err = InsertNonDicDB(UserNew, structFromReq, db, params)
				if !err {
					var r UserInCommon
					err = db.C("users").Find(bson.M{
						"_id": newId}).Select(SelectUserInCommon).One(&r)
					if !err {
						v.FieldByName("User").Set(reflect.ValueOf(r))
						// Proceed to tokens
						var r1 TokensInCommon
						r1, err = SetDeviceTokens(db, structFromReq, params)
						if !err {
							err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
						}
					}
				}
			}
		case signIn:
			var r UserInCommon
			err = db.C("users").Find(criteria).Select(SelectUserInCommon).One(&r)
			if !err {
				v.FieldByName("User").Set(reflect.ValueOf(r))
				// Proceed to tokens
				var r1 TokensInCommon
				r1, err = SetDeviceTokens(db, structFromReq, params)
				if !err {
					err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
				}
			}
		case renewTokens:
			var r TokensInCommon
			r, err = SetDeviceTokens(db, structFromReq, params)
			if !err {
				err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r))
				if !err {
					err = SetResBodyPart(v.FieldByName("UserId"), "UserId", reflect.ValueOf(params["user_id"]))
				}
			}
		case newDeviceInfo:
			// The only case a new deviceInfo created by client is after signing in.
			var newId bson.ObjectId
			newId, err = InsertNonDicDB(DeviceInfoNew, structFromReq, db, params)
			if !err {
				var r DeviceInfoInCommon
				err = db.C("deviceInfos").Find(bson.M{
					"_id": newId}).Select(SelectDeviceInfoInCommon).One(&r)
				if !err {
					err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
				}
			}
		case oneDeviceInfoSortOption:
			// When log in on a new device or a device that the account has been deleted before, both indicate no deviceInfo on the device, 
			var n int
			n, err = db.C("deviceInfos").Find(criteria).Count()
			if n <= 0 || err {
				var selector bson.M
				selector, err = UpdateNonDicDB(DeviceInfoUpdateSortOption, structFromReq, db, params)
				if !err {
					var r DeviceInfoInCommon
					err = db.C("deviceInfos").Find(selector).Select(SelectDeviceInfoInCommon).One(&r)
					if !err {
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
					}
				}
			}
		case oneDeviceInfoLang:
			// When log in on a new device or a device that the account has been deleted before, both indicate no deviceInfo on the device, 
			var n int
			n, err = db.C("deviceInfos").Find(criteria).Count()
			if n <= 0 || err {
				var selector bson.M
				selector, err = UpdateNonDicDB(DeviceInfoUpdateLang, structFromReq, db, params)
				if !err {
					var r DeviceInfoInCommon
					err = db.C("deviceInfos").Find(selector).Select(SelectDeviceInfoInCommon).One(&r)
					if !err {
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
					}
				}
			}
		// Resetting password goes through html, not here.
		// Activating only changes email after user being activated. It goes through html, not here, either.
		case sync:
			userId := bson.ObjectIdHex(params["user_id"])
			// Sync User
			myUser := UserInCommon
			err = db.C("users").Find(bson.M{
				"_id": userId}).Select(SelectUserInCommon).One(&myUser)
			if !err {
				if myUser.VersionNo > x.FieldByName("User").FieldByName("VersionNo").Int() {
					// No need to send userInCommon to client in other cases.
					err = SetResBodyPart(v.FieldByName("User"), "User", reflect.ValueOf(myUser))
					if !err {
						// Sync DeviceInfo
						myDeviceInfo := DeviceInfoInCommon
						err = db.C("deviceInfos").Find(bson.M{
							"belongTo": userId,
							"deviceUUID": x.FieldByName("DeviceUUID").String(),
							"_id": x.FieldByName("DeviceInfo").FieldByName("_id")}).Select(SelectDeviceInfoInCommon).One(&myDeviceInfo)
						if !err {
							err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(myDeviceInfo))
							if !err {
								var cardListDB []CardsVerList
								cardListDB, err = GetDbCardVerList(db, userId)
								if !err {
									err = GetCardListDifference(db, cardListDB, x.FieldByName("CardList"), v.FieldByName("CardList"), v.FieldByName("CardToDelete"))
								}
							}
						}
					}
				}
			}
		case newCard:
			var r []CardInCommon
			var r1 []CardInCommon
			r, r1, err = InsertNewCard(db, structFromReq, params)
			if len(r) > 0 {
				// Card with same content found, return this card to client.
				err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
			}
		case oneCard:
			var t CardInCommon
			var r []CardInCommon
			err = db.C("cards").Find(criteria).One(&t)
			if err == mgo.ErrNotFound {
				// No card with this exists. It indicates the card sent from client has been deleted already. Assign a new id, if it's an update operation.
				// Here treat it as a new card and overwrite it on client.
				var r1 []CardInCommon
				r, r1, err = InsertNewCard(db, structFromReq, params)
				if len(r) > 0 {
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			} else if !err {
				// Update the card according to the conflict-resolving rule.
				decisionCode := HandleCardVersionConflict("update", x.FieldByName("Card"), t)
				if decisionCode == ConflictCreateAnotherInDB {
					// Check unique first and every step like insert a new card.
					// Need to overwrite the record on client with this id as well.
					var r1 []CardInCommon
					r, r1, err = InsertNewCard(db, structFromReq, params)
					if len(r) > 0 {
						r = append(r, t)
						err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
					} else if len(r1) > 0 {
						// Updated content is duplicated with other existing records in db. Deliver those to client. 
						r1 = append(r1, t)
						// Clear potential err in previoue steps.
						err = nil
						err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r1))
					}
				} else if decisionCode == NoConflictOverwriteDB {
					// No need to check uniqueness here. If there is duplicate card in db, just let the user do what he wants.
					var selector bson.M
					selector, err = UpdateNonDicDB(CardUpdate, structFromReq, db, params)
					if !err {
						err = db.C("cards").Find(selector).Select(SelectCardInCommon).All(&r)
						if !err {
							err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
						}
					}
				} else if decisionCode == ConflictOverwriteClient {
					r = make([]CardInCommon, 0)
					r = append(r, t)
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			}	
		case dict:
			// TBD
		default:
			e = error.New("Request not recogniized.")
	} else if m == "DELETE" {
		// One case only: DELETE /users/:user_id/cards/:card_id
		if route == oneCard {
			decision = HandleVersionConflict("delete", &x.FieldByName("Card"), &v)
		} else {
			decision = MethodNotAllowed
		}
	}
	return
}

func GetDbCardVerList(db *mgo.Database, userId bson.ObjectId) (r []CardsVerList, err error) {
	// Shoud identify the err returned for further possibities.
	err = db.C("Cards").Find(bson.M{
		"belongTo": userId,
		"isDeleted": false}).Select(bson.M{
		"_id": 1,
		"versionNo": 1}).All(&r)
	return
}

func GetCardListDifference(db *mgo.Database, cardListDB []CardsVerList, cardListReq reflect.Value, cardListRes reflect.Value, cardsToDelete reflect.Value) (errs error) {
	// cardListReq: []CardsVerList, cardListRes: []CardInCommon, cardsToDelete: []bson.ObjectId
	// The difference is to overwrite/add docs on client
	// Only non-deleted cards on db here, see PreprocessRequest.
	if !cardListRes.CanSet() || !cardsToDelete.CanSet() {
		return
	}
	for x := range cardListDB {
		// Reset the flag
		hasMatch = false
		// Search if the specific id exists 
		for y := 0; y < cardListReq.Len(); y++ {
			if x.Id == bson.ObjectIdHex(cardListReq.Index(y).FielfByName("Id").String()) {
				hasMatch = true
				if x.VersionNo == cardListReq.Index(y).FielfByName("VersionNo").Int() {
					// Same version, do nothing
				} else {
					// Overwrite the doc on client even if x.VersionNo < cardListReq.Index(y).FielfByName("VersionNo").Int() which indicates something wrong on client.
					err := PushCardToCardList(db, cardListRes, x.Id)
					if err != nil {
						return
					}
				}
				break
			}
		}
		// If the specific id not exists, push to list
		if hasMatch == false {
			err := PushCardToCardList(db, cardListRes, x.Id)
			if err != nil {
				return
			}
		}
	}
	// Just find out which ones are not on server
	for k := 0; k < cardListReq.Len(); k++ {
		// Search if the specific id exists 
		for j := range cardListDB {
			if bson.ObjectIdHex(cardListReq.Index(k).FielfByName("Id").String()) == j.Id {
				break
			}
		}
		// If the specific id not exists, push to delete list anyway
		cardsToDelete = reflect.Append(cardsToDelete, cardListReq.Index(k).FielfByName("Id"))
	}
	return
}

func PushCardToCardList(db *mgo.Database, cardList reflect.Value, cardID bson.ObjectId) (err error) {
	// cardList: []CardInCommon
	var aCard CardInCommon
	err = db.C("cards").Find(bson.M{
		"_id": cardID,
		"isDeleted": false}).Select(bson.M{
			"isDeleted": 0}).One(&aCard)
	if !err {
		cardList = reflect.Append(cardList, reflect.ValueOf(aCard))
	}
	// err == mgo.ErrNotFound is still an error here.
	return
}

func SetResBodyPart(partToSet reflect.Value, fieldNameToSet string, valueIn interface{}) (err error) {
    if partToSet.CanSet() {
		v := reflect.ValueOf(valueIn)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		v.FieldByName(fieldNameToSet).Set(v)
	} else {
		msg := string.Join([]string{"Structure for response not able to be set:", fieldNameToSet}, " ")
		err = errors.New(msg)
	}
    return
}

func EncodeAndResponse(structInRes interface{}, rw *martini.ResponseWriter) {
	body, err := json.Marshal(structInRes)
	if err != nil {
		http.Error(rw, "Encoding response body failed", StatusInternalServerError)
	} else {
		rw.Write(body)
	}
}

func InsertNewCard(db *mgo.Database, structFromReq interface{}, params *martini.Params) (inserted []CardInCommon, duplicated []CardInCommon, err error) {
	x := reflect.ValueOf(structFromReq)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	var isUnique bool
	isUnique, duplicated, err = CheckCardUniqueness(db, params, x.FieldByName("Card"))
	if isUnique {
		// New card is unique, insert it into db. err == mgo.ErrNotFound
		err = nil
		var newId bson.ObjectId
		newId, err = InsertNonDicDB(CardNew, structFromReq, db, params)
		if !err {
			err = db.C("cards").Find(bson.M{
				"_id": newId}).Select(SelectCardInCommon).All(&inserted)
		}
	}
}
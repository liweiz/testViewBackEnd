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
		case oneDeviceInfo:
			err = db.C("deviceInfos").Find(criteria).Count()
		
		// Resetting password goes through html, not here.
		// Activating only changes email after user being activated. It goes through html, not here, either.
		case sync:
			userId := bson.ObjectIdHex(params["user_id"])
			// Sync User
			myUser := UserInCommon
			errUser := db.C("users").Find(bson.M{
				"_id": userId}).Select(SelectUserInCommon).One(&myUser)
			outUser := HandleUserConflict(x.FieldByName("User").FieldByName("VersionNo").Int(), myUser)
			if outUser == OverWriteClient {
				if v.FieldByName("User").CanSet() {
					v.FieldByName("User").Set(reflect.ValueOf(myUser))
				} else {
					errUser1 := db.C("users").Update()
				}
			}
			// Sync DeviceInfo
			myDeviceInfo := DeviceInfoInCommon
			errDeviceInfo := db.C("deviceinfos").Find(bson.M{
				"belongTo": params["user_id"],
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"_id": x.FieldByName("DeviceInfo").FieldByName("_id")}).Select(SelectDeviceInfoInCommon).One(&myDeviceInfo)
			outDeviceInfo := HandleDeviceInfoConflict(x.FieldByName("DeviceInfo").FieldByName("VersionNo").Int(), myDeviceInfo)
			if outDeviceInfo == OverWriteClient {
				if v.FieldByName("DeviceInfo").CanSet() {
					v.FieldByName("DeviceInfo").Set(reflect.ValueOf(myDeviceInfo))
				} else {
					decision = InternalServerError
				}
			}
			// Sync Card
			cardListDB, errCard := GetDbCardVerList(db, userId)			
			if errCard != nil && errCard != mgo.ErrNotFound {
				decision = InternalServerError
			} else {
				// For CardList
				errCard1 := GetCardListDifference(db, true, cardListDB, &x.FieldByName("CardList"), &v.FieldByName("CardList"), nil)
				if errCard1 == mgo.ErrNotFound {
					decision = InternalServerError
				} else if errCard1 != nil {
					// No cards on server. No need to set the CardList in response structure since it's originally nil.
				} else {
					errCard2 := GetCardListDifference(db, false, cardListDB, &x.FieldByName("CardList"), nil, &v.FieldByName("CardToDelete"))
					if errCard2 != nil {
						// ****No such card found
						decision = InternalServerError
					} else {
						decision = OKAndSync
					}
				}
			}
		case newCard:
			// Check uniqueness first
			if CheckCardUniqueness(db, params, x.FieldByName("Card")) {
				// New card is unique, insert it into db.
				docToSave := PrepareDocForDB(CardNew, x.FieldByName("Card"), nil, params)
				dts := reflect.ValueOf(docToSave)
				if dts.Kind() == reflect.Ptr {
					dts = dts.Elem()
				}
				err := db.C("cards").Insert()
			} else {
				decision = ConflictCardAlreadyExists
			}
		case oneCard:
			// var aCard Card
			err := db.C("cards").Find(criteria).One(&v)
			if err != nil {
				decision = NotFoundOnly
			} else {
				// Check uniqueness first
				CheckDocUniqueness(db, )
				// Create
				docToSave := PrepareDocForDB(CardNew, x.Card, nil, params)
				dts := reflect.ValueOf(docToSave)
				if dts.Kind() == reflect.Ptr {
					dts = dts.Elem()
				}
				if dts.Kind() == reflect.Map {
					err1 := db.C("users").Insert(docToSave)
					if err1 != nil {
						decision = InternalServerError
					} else {
						decision = OkAndUser
					}
				} else {
					decision = InternalServerError
				}
				decision = HandleVersionConflict("update", &x.FieldByName("Card"), &v)
			}
		// No need here: case allCards:	
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

func GetCardListDifference(db *mgo.Database, baseIsClient bool, cardListDB []CardsVerList, cardListReq []CardsVerList, cardListRes []*CardInCommon, cardsToDelete []bson.ObjectId) (errs []error) {
	var hasMatch bool
	if baseIsClient == true {
		// The difference is to overwrite/add docs on client
		// Only non-deleted cards on db here, see PreprocessRequest.
		for x := range cardListDB {
			// Reset the flag
			hasMatch = false
			// Search if the specific id exists 
			for y := range cardListReq {
				if x.Id == y.Id {
					hasMatch = true
					if x.VersionNo == y.VersionNo {
						// Same version, do nothing
					} else {
						// Overwrite the doc on client
						err := PushCardToCardList(db, cardListRes, x.Id)
						if err != nil {
							if errs == nil {
								errs = make([]error, 0)
							}
							errs = append(errs, err)
						}
					}
					break
				}
			}
			// If the specific id not exists, push to list
			if hasMatch == false {
				err := PushCardToCardList(db, cardListRes, x.Id)
				if err != nil {
					if errs == nil {
						errs = make([]error, 0)
					}
					errs = append(errs, err)
				}
			}
		}
	} else {
		// Just find out which ones are not on server
		for x := range cardListReq {
			// Search if the specific id exists 
			for y := range cardListDB {
				if x.Id == y.Id {
					break
				}
			}
			// If the specific id not exists, push to delete list anyway
			cardsToDelete = append(cardsToDelete, x)
		}
	}
	return
}

func PushCardToCardList(db *mgo.Database, cardList []*CardInCommon, cardID bson.ObjectId) (err error) {
	var aCard CardInCommon
	err = db.C("cards").Find(bson.M{
		"_id": cardID}).Select(bson.M{
			"isDeleted": 0}).One(&aCard)
	if !err {
		cardList = append(cardList, &aCard)
	}
	// Need to make decision outside of this func according to err == mgo.ErrNotFound
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
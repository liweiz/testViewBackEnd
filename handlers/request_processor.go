package testView

import (
	"encoding/json"
	"errors"
	"github.com/codegangsta/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func ProcessedResponseGenerator() martini.Handler {
	return func(db *mgo.Database, route int, v *Vehicle, req *http.Request, params martini.Params, rw martini.ResponseWriter, logger *log.Logger) {
		err := ProcessRequest(db, route, v.Criteria, v.ReqStruct, req, v.ResStruct, params)
		var s string
		if err == nil {
			// Add successfully processed request to RequestProcessed.
			_, err = InsertNonDicDB(RequestProcessedNew, v.ReqStruct, db, params)
			if err != nil {
				s = strings.Join([]string{"Failed to insert RequestProcessed, but request has been successfully processed by server. Sent 200 to client so that no more same request needs to be sent.", err.Error()}, "=> ")
				// Only log it, no need to send to client.
				WriteLog(s, logger)
				// Reset err and s.
				err = nil
				s = ""
			}
			// Send response.
			rw.Header().Set("Content-Type", "application/json")
			var j []byte
			j, err = json.Marshal(v.ReqStruct)
			if err == nil {
				// Response size, any usage???
				_, err = rw.Write(j)
			}
			if err != nil {
				s = strings.Join([]string{"Failed to generate response, but request has been successfully processed by server.", err.Error()}, "=> ")
			}
		}
		if err != nil {
			if s == "" {
				WriteLog(err.Error(), logger)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			} else {
				WriteLog(s, logger)
				http.Error(rw, s, http.StatusServiceUnavailable)
			}
		}
	}
}

func ProcessRequest(db *mgo.Database, route int, criteria bson.M, structFromReq interface{}, req *http.Request, structForRes interface{}, params martini.Params) (err error) {
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
		switch route {
		case oneCard:
			var resultCard CardInCommon
			err = db.C("cards").Find(criteria).Select(GetSelector(SelectCardInCommon)).One(&resultCard)
			if err == nil {
				err = SetResBodyPart(v, "Cards", []CardInCommon{resultCard})
			}
		case oneUser:
			var resultUser UserInCommon
			err = db.C("users").Find(criteria).Select(GetSelector(SelectUserInCommon)).One(&resultUser)
			if err == nil {
				err = SetResBodyPart(v, "User", resultUser)
			}
		case oneDeviceInfo:
			var resultDeviceInfo DeviceInfoInCommon
			err := db.C("deviceInfos").Find(bson.M{
				"belongTo": params["user_id"],
				// /users/:user_id/devices/:device_id
				"_id": params["device_id"]}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&resultDeviceInfo)
			if err == nil {
				err = SetResBodyPart(v, "DeviceInfo", resultDeviceInfo)
			}
		// case dicTranslation:
		// case dicDetail:
		default:
			err = errors.New("Request not recogniized.")
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
			if n > 0 && err == nil {
				var newId bson.ObjectId
				newId, err = InsertNonDicDB(UserNew, structFromReq, db, params)
				if err == nil {
					var r UserInCommon
					err = db.C("users").Find(bson.M{
						"_id": newId}).Select(GetSelector(SelectUserInCommon)).One(&r)
					if err == nil {
						v.FieldByName("User").Set(reflect.ValueOf(r))
						// Proceed to tokens
						var r1 TokensInCommon
						r1, err = SetDeviceTokens(db, structFromReq, params)
						if err == nil {
							err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
						}
					}
				}
			}
		case signIn:
			var r UserInCommon
			err = db.C("users").Find(criteria).Select(GetSelector(SelectUserInCommon)).One(&r)
			if err == nil {
				v.FieldByName("User").Set(reflect.ValueOf(r))
				// Proceed to tokens
				var r1 TokensInCommon
				r1, err = SetDeviceTokens(db, structFromReq, params)
				if err == nil {
					err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
				}
			}
		case renewTokens:
			var r TokensInCommon
			r, err = SetDeviceTokens(db, structFromReq, params)
			if err == nil {
				err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r))
				if err == nil {
					err = SetResBodyPart(v.FieldByName("UserId"), "UserId", reflect.ValueOf(params["user_id"]))
				}
			}
		case newDeviceInfo:
			// The only case a new deviceInfo created by client is after signing in.
			var newId bson.ObjectId
			newId, err = InsertNonDicDB(DeviceInfoNew, structFromReq, db, params)
			if err == nil {
				var r DeviceInfoInCommon
				err = db.C("deviceInfos").Find(bson.M{
					"_id": newId}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&r)
				if err == nil {
					err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
				}
			}
		case oneDeviceInfoSortOption:
			// When log in on a new device or a device that the account has been deleted before, both indicate no deviceInfo on the device,
			var n int
			n, err = db.C("deviceInfos").Find(criteria).Count()
			if n <= 0 || err == mgo.ErrNotFound {
				var selector bson.M
				selector, err = UpdateNonDicDB(DeviceInfoUpdateSortOption, structFromReq, db, params)
				if err == nil {
					var r DeviceInfoInCommon
					err = db.C("deviceInfos").Find(selector).Select(GetSelector(SelectDeviceInfoInCommon)).One(&r)
					if err == nil {
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
					}
				}
			}
		case oneDeviceInfoLang:
			// When log in on a new device or a device that the account has been deleted before, both indicate no deviceInfo on the device,
			var n int
			n, err = db.C("deviceInfos").Find(criteria).Count()
			if n <= 0 || err == mgo.ErrNotFound {
				var selector bson.M
				selector, err = UpdateNonDicDB(DeviceInfoUpdateLang, structFromReq, db, params)
				if err == nil {
					var r DeviceInfoInCommon
					err = db.C("deviceInfos").Find(selector).Select(GetSelector(SelectDeviceInfoInCommon)).One(&r)
					if err == nil {
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
					}
				}
			}
		// Resetting password goes through html, not here.
		// Activating only changes email after user being activated. It goes through html, not here, either.
		case sync:
			userId := bson.ObjectIdHex(params["user_id"])
			// Sync User
			var myUser UserInCommon
			err = db.C("users").Find(bson.M{
				"_id":       userId,
				"isDeleted": false}).Select(GetSelector(SelectUserInCommon)).One(&myUser)
			if err == nil {
				if myUser.VersionNo > x.FieldByName("User").FieldByName("VersionNo").Int() {
					// No need to send userInCommon to client in other cases.
					err = SetResBodyPart(v.FieldByName("User"), "User", reflect.ValueOf(myUser))
					if err == nil {
						// Sync DeviceInfo
						var myDeviceInfo DeviceInfoInCommon
						err = db.C("deviceInfos").Find(bson.M{
							"belongTo":   userId,
							"deviceUUID": x.FieldByName("DeviceUUID").String(),
							"_id":        x.FieldByName("DeviceInfo").FieldByName("_id")}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&myDeviceInfo)
						if err == nil {
							err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(myDeviceInfo))
							if err == nil {
								var cardListDB []CardsVerList
								cardListDB, err = GetDbCardVerList(db, userId)
								if err == nil {
									err = GetCardListDifference(db, cardListDB, x.FieldByName("CardList"), v.FieldByName("CardList"), v.FieldByName("CardToDelete"))
								}
							}
						}
					}
				}
			}
		case newCard:
			var r []CardInCommon
			r, _, err = InsertNewCard(db, structFromReq, params)
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
				r, _, err = InsertNewCard(db, structFromReq, params)
				if len(r) > 0 {
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			} else if err == nil {
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
					if err == nil {
						err = db.C("cards").Find(selector).Select(GetSelector(SelectCardInCommon)).All(&r)
						if err == nil {
							err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
						}
					}
				} else if decisionCode == ConflictOverwriteClient {
					r = make([]CardInCommon, 0)
					r = append(r, t)
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			}
		//case dict:
		// For words only. Client posts WordsText to exchange its id and translation/target list.
		// var r []DicTextInCommon
		// err = db.C("cards").Find(criteria).Select(SelectDicInCommon).All(&r)
		// if err == nil {

		// }
		default:
			err = errors.New("Request not recogniized.")
		}
	} else if m == "DELETE" {
		// One case only: DELETE /users/:user_id/cards/:card_id
		if route == oneCard {
			var t CardInCommon
			var r []CardInCommon
			err = db.C("cards").Find(criteria).One(&t)
			if err == mgo.ErrNotFound {
				// Return ok since no such card exists among non-deleted ones.
			} else if err == nil {
				decisionCode := HandleCardVersionConflict("delete", x.FieldByName("Card"), t)
				if decisionCode == NoConflictOverwriteDB {
					err = db.C("cards").Remove(criteria)
				} else if decisionCode == ConflictOverwriteClient {
					r = make([]CardInCommon, 0)
					r = append(r, t)
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			}
		} else {
			err = errors.New("Request not recogniized.")
		}
	}
	return
}

func GetDbCardVerList(db *mgo.Database, userId bson.ObjectId) (r []CardsVerList, err error) {
	// Shoud identify the err returned for further possibities.
	err = db.C("Cards").Find(bson.M{
		"belongTo":  userId,
		"isDeleted": false}).Select(bson.M{
		"_id":       1,
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
		hasMatch := false
		xx := cardListDB[x]
		// Search if the specific id exists
		for y := 0; y < cardListReq.Len(); y++ {
			if xx.Id == bson.ObjectIdHex(cardListReq.Index(y).FieldByName("Id").String()) {
				hasMatch = true
				if xx.VersionNo == cardListReq.Index(y).FieldByName("VersionNo").Int() {
					// Same version, do nothing
				} else {
					// Overwrite the doc on client even if x.VersionNo < cardListReq.Index(y).FieldByName("VersionNo").Int() which indicates something wrong on client.
					err := PushCardToCardList(db, cardListRes, xx.Id)
					if err != nil {
						return
					}
				}
				break
			}
		}
		// If the specific id not exists, push to list
		if hasMatch == false {
			err := PushCardToCardList(db, cardListRes, xx.Id)
			if err != nil {
				return
			}
		}
	}
	// Just find out which ones are not on server
	for k := 0; k < cardListReq.Len(); k++ {
		// Search if the specific id exists
		for j := range cardListDB {
			jj := cardListDB[j]
			if bson.ObjectIdHex(cardListReq.Index(k).FieldByName("Id").String()) == jj.Id {
				break
			}
		}
		// If the specific id not exists, push to delete list anyway
		cardsToDelete = reflect.Append(cardsToDelete, cardListReq.Index(k).FieldByName("Id"))
	}
	return
}

func PushCardToCardList(db *mgo.Database, cardList reflect.Value, cardID bson.ObjectId) (err error) {
	// cardList: []CardInCommon
	var aCard CardInCommon
	err = db.C("cards").Find(bson.M{
		"_id":       cardID,
		"isDeleted": false}).Select(bson.M{
		"isDeleted": 0}).One(&aCard)
	if err == nil {
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
		msg := strings.Join([]string{"Structure for response not able to be set:", fieldNameToSet}, " ")
		err = errors.New(msg)
	}
	return
}

func InsertNewCard(db *mgo.Database, structFromReq interface{}, params martini.Params) (inserted []CardInCommon, duplicated []CardInCommon, err error) {
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
		if err == nil {
			err = db.C("cards").Find(bson.M{
				"_id": newId}).Select(GetSelector(SelectCardInCommon)).All(&inserted)
		}
	}
	return
}

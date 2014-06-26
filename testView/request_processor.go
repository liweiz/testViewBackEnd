package testView

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	// "github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	// "path"
	"os"
	"reflect"
	"strings"
	"time"
)

func ProcessedResponseGenerator(route int, setRequestProcessed bool) martini.Handler {
	return func(db *mgo.Database, v *Vehicle, req *http.Request, params martini.Params, rw http.ResponseWriter, logger *log.Logger) {
		err := ProcessRequest(db, route, v.Criteria, v.ReqStruct, req, v.ResStruct, params)
		var s string
		if err == nil {
			// Add successfully processed request to RequestProcessed, if needed.
			if setRequestProcessed {
				_, err = InsertNonDicDB(RequestProcessedNew, v.ReqStruct, db, bson.ObjectIdHex(params["user_id"]))
				if err != nil {
					s = strings.Join([]string{"Failed to insert RequestProcessed, but request has been successfully processed by server. Sent 200 to client so that no more same request needs to be sent.", err.Error()}, "=> ")
					// Only log it, no need to send to client.
					WriteLog(s, logger)
					// Reset err and s.
					err = nil
					s = ""
				}
			}
			// Send response.
			rw.Header().Set("Content-Type", "application/json")
			var j []byte
			j, err = json.Marshal(v.ResStruct)
			if err == nil {
				// Response size, any usage???
				_, err = rw.Write(j)
				fmt.Println("code:", 200)
				os.Stdout.Write(j)
			}
			if err != nil {
				s = strings.Join([]string{"Failed to generate response, but request has been successfully processed by server.", err.Error()}, "=> ")
			}
		}
		if err != nil {
			if s == "" {
				WriteLog(err.Error(), logger)
				fmt.Println("err to log: ", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			} else {
				WriteLog(s, logger)
				http.Error(rw, s, http.StatusServiceUnavailable)
			}
		}
	}
}

func ProcessedResponseGeneratorPasswordResetting() martini.Handler {
	return func(db *mgo.Database, v *Vehicle, req *http.Request, params martini.Params, rw http.ResponseWriter, logger *log.Logger) {
		err := ProcessRequest(db, PasswordResetting, v.Criteria, v.ReqStruct, req, nil, params)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
		} else if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
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
		case OneCard:
			var resultCard CardInCommon
			err = db.C("cards").Find(criteria).Select(GetSelector(SelectCardInCommon)).One(&resultCard)
			if err == nil {
				err = SetResBodyPart(v, "Cards", reflect.ValueOf([]CardInCommon{resultCard}))
			}
		case OneDeviceInfo:
			var resultDeviceInfo DeviceInfoInCommon
			err := db.C("deviceInfos").Find(bson.M{
				"belongTo": bson.ObjectIdHex(params["user_id"]),
				// /users/:user_id/devices/:device_id
				"_id": bson.ObjectIdHex(params["device_id"])}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&resultDeviceInfo)
			if err == nil {
				err = SetResBodyPart(v, "DeviceInfo", reflect.ValueOf(resultDeviceInfo))
			}
		case OneUser:
			var resultUser UserInCommon
			resStruct = ResUser{}
			err = db.C("users").Find(bson.M{
				"_id": bson.ObjectIdHex(params["user_id"])}).Select(GetSelector(SelectUserInCommon)).One(&resultUser)
			if err == nil {
				resStruct.User = resultUser
				err = SetResBodyPart(reflect.ValueOf(resStruct), "User", reflect.ValueOf(resultUser))
			}
		// case dicTranslation:
		// case dicDetail:
		default:
			err = errors.New("Request not recognized.")
		}
	} else if m == "POST" {
		switch route {
		/*
			SignUp flow:
			1. Client sends email and password to server
			2. Create user in db if everything's ok. Otherwise, return error message
			3.1 Client does not receive response, while server has successfully processed the request. User will have to sign in.
			3.2 Client receives the user, lets the user set languagePair and sends it as deviceInfo to server.
			4.1 The deviceInfo created on client is not processed successfully by server (such as request not received), follow steps in deviceInfo initialization flow on each client.
			4.2 Server creates deviceInfo in db if everything's ok. Otherwise, return error message.

			SignIn flow:
			1. Client sends email and password to server
			2. Get user in db if everything's ok. Otherwise, return error message.
			2.1 No previous data on client found. Start sync process. If no deviceInfo on client, use the default one on server.
			2.2 Previous data on client. Start sync process.

			DeviceInfo initialization flow on each client
			This is for getting the deviceInfo the first time on each client for a given account. It's either triggered by user creating a new one on client and posting it to server or nothing on client and a sync request sent to get it.
			record below is not necessarily successfully post to server.
			1. no record on client, no record on server. Let user create one on client.
			2. no record on client, record exists on server. Get the default record on server and deliver to client.
			3. record exists on client, no record on server. Post the record on client to server.
			4. record exists on both client and server. Always overwrite the one on server since deviceInfo is device and account specified. In this case, record on client is modified from the one on server.

			So a new user's deviceInfo is created on client by user and stored to db. A existing user's deviceInfo is delivered by sync process or providing the list for client to choose. When there is only one existing deviceInfo on server, send this to client.
		*/
		// case SignUp:
		// 	var aUser User
		// 	// Check if already in use
		// 	err = db.C("users").Find(criteria).One(&aUser)
		// 	if err == mgo.ErrNotFound {
		// 		fmt.Println("existing user not found")
		// 		// Reset err = nil
		// 		err = nil
		// 		// Create a new user
		// 		var newId bson.ObjectId
		// 		newId, err = InsertNonDicDB(UserNew, structFromReq, db, "")
		// 		if err == nil {
		// 			// Get and put the new user to response body.
		// 			var r UserInCommon
		// 			err = db.C("users").Find(bson.M{
		// 				"_id": newId}).Select(GetSelector(SelectUserInCommon)).One(&r)
		// 			if err == nil {
		// 				err = SetResBodyPart(v.FieldByName("User"), "User", reflect.ValueOf(r))
		// 				if err == nil {
		// 					// Proceed to tokens
		// 					var r1 TokensInCommon
		// 					r1, err = SetGetDeviceTokens(r.Id, structFromReq, db)
		// 					if err == nil {
		// 						err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
		// 					}
		// 				}
		// 			}
		// 		}
		// 	} else if err == nil {
		// 		err = errors.New("User already exists.")
		// 	}
		// case SignIn:
		// 	var r UserInCommon
		// 	err = db.C("users").Find(criteria).Select(GetSelector(SelectUserInCommon)).One(&r)
		// 	if err == nil {
		// 		err = SetResBodyPart(v.FieldByName("User"), "User", reflect.ValueOf(r))
		// 		if err == nil {
		// 			// Proceed to tokens, everytime signIn from a client, tokens have to be refreshed.
		// 			var r1 TokensInCommon
		// 			r1, err = SetGetDeviceTokens(r.Id, structFromReq, db)
		// 			if err == nil {
		// 				err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r1))
		// 			}
		// 		}
		// 	}
		case ForgotPassword:
			var aUser User
			err = db.C("users").Find(criteria).One(&aUser)
			if err == mgo.ErrNotFound {
				err = errors.New("User does not exist.")
			} else if err == nil {
				err = SendEmail(EmailForPasswordResetting, db, req, params, aUser.Id)
			}
		case RenewTokens:
			var r TokensInCommon
			r, err = SetGetDeviceTokens(bson.ObjectIdHex(params["user_id"]), structFromReq, db)
			if err == nil {
				err = SetResBodyPart(v.FieldByName("Tokens"), "Tokens", reflect.ValueOf(r))
				if err == nil {
					err = SetResBodyPart(v.FieldByName("UserId"), "UserId", reflect.ValueOf(params["user_id"]))
				}
			}
		case NewDeviceInfo:
			// The only case a new deviceInfo created by client is after signing in.
			userId := bson.ObjectIdHex(params["user_id"])
			var deviceInfoIdHere bson.ObjectId
			var dd DeviceInfo
			dd, err = EnsureDeviceInfoUniqueness(db, userId, x.FieldByName("DeviceUUID").String())
			if err == nil {
				// DeviceInfo already exists, please update it if needed.
				err = db.C("deviceInfos").UpdateId(dd.Id, bson.M{
					"$set": bson.M{
						"lastModified": time.Now().UnixNano(),
						"sortOption":   x.FieldByName("DeviceInfo").FieldByName("SortOption").String(),
						"sourceLang":   x.FieldByName("DeviceInfo").FieldByName("SourceLang").String(),
						"targetLang":   x.FieldByName("DeviceInfo").FieldByName("TargetLang").String()}})
				if err == nil {
					deviceInfoIdHere = dd.Id
				}
			} else if err == mgo.ErrNotFound {
				var newId bson.ObjectId
				newId, err = InsertNonDicDB(DeviceInfoNew, structFromReq, db, userId)
				if err == nil {
					deviceInfoIdHere = newId
				}
			}
			if err == nil {
				var r DeviceInfoInCommon
				err = db.C("deviceInfos").Find(bson.M{
					"_id": deviceInfoIdHere}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&r)
				fmt.Println("DeviceInfoInCommon: ", r)
				if err == nil {
					err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
				}
			}
		case OneDeviceInfo:
			// When log in on a new device or a device that the account has been deleted before, both indicate no deviceInfo on the device,
			var q DeviceInfo
			err = db.C("deviceInfos").Find(criteria).One(&q)
			if err == nil {
				var selector bson.M
				selector, err = UpdateNonDicDB(DeviceInfoUpdate, structFromReq, db, params, bson.ObjectIdHex(params["user_id"]))
				if err == nil {
					var r DeviceInfoInCommon
					err = db.C("deviceInfos").Find(selector).Select(GetSelector(SelectDeviceInfoInCommon)).One(&r)
					if err == nil {
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(r))
					}
				}
			}
		case PasswordResetting:
			userId := bson.ObjectIdHex(params["user_id"])
			_, err = UpdateNonDicDB(UserUpdatePassword, structFromReq, db, params, userId)
			if err == nil {
				_, _ = UpdateNonDicDB(UserRemovePasswordUrlCodeRoutinely, nil, db, params, userId)
				_ = RemoveUserPasswordUrlCodeOne(db, params, userId)
			}
		case NewCard:
			var r []CardInCommon
			r, _, err = InsertNewCard(db, structFromReq, params)
			fmt.Println("number of new card in response: ", len(r))
			if err == nil {
				// Card inserted as new card, return this card to client.
				err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				_ = ScanToFindAndUpdateEmptyTextId(db, bson.ObjectIdHex(params["user_id"]))
			}
		case OneCard:
			var t CardInCommon
			var r []CardInCommon
			err = db.C("cards").Find(criteria).Select(GetSelector(SelectCardInCommon)).One(&t)
			if err == mgo.ErrNotFound {
				// No card with this exists. It indicates the card sent from client has been deleted already. Assign a new id and treat it as a new card operation, if it's an update operation. Overwrite it on client.
				r, _, err = InsertNewCard(db, structFromReq, params)
				if err == nil {
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
					_ = ScanToFindAndUpdateEmptyTextId(db, bson.ObjectIdHex(params["user_id"]))
				}
			} else if err == nil {
				var isUnique bool
				isUnique, _, err = CheckCardUniqueness(db, params, x.FieldByName("Card"))
				if isUnique {
					// Update the card according to the conflict-resolving rule.
					decisionCode := HandleCardVersionConflict("update", structFromReq, t)
					if decisionCode == ConflictCreateAnotherInDB {
						// Check unique first and every step like insert a new card.
						// Need to overwrite the record on client with this id as well.
						var newId bson.ObjectId
						newId, err = InsertNonDicDB(CardNew, structFromReq, db, bson.ObjectIdHex(params["user_id"]))
						if err == nil {
							err = db.C("cards").Find(bson.M{
								"_id": newId}).Select(GetSelector(SelectCardInCommon)).All(&r)
							if err == nil {
								r = append(r, t)
								err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
							}
						}
					} else if decisionCode == NoConflictOverwriteDB {
						// No need to check uniqueness here. If there is duplicate card in db, just let the user do what he wants.
						var selector bson.M
						selector, err = UpdateNonDicDB(CardUpdate, structFromReq, db, params, bson.ObjectIdHex(params["user_id"]))
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
					_ = ScanToFindAndUpdateEmptyTextId(db, bson.ObjectIdHex(params["user_id"]))
				}
			}
			// Activating only changes email after user being activated. It goes through html, not here, either.
		case Sync:
			userId := bson.ObjectIdHex(params["user_id"])
			// Sync User
			var myUser UserInCommon
			err = db.C("users").Find(bson.M{
				"_id":       userId,
				"isDeleted": false}).Select(GetSelector(SelectUserInCommon)).One(&myUser)
			if err == nil {
				// Overwrite client's with server's anyway
				err = SetResBodyPart(v.FieldByName("User"), "User", reflect.ValueOf(myUser))
				if err == nil {
					// Sync DeviceInfo
					var myDeviceInfo DeviceInfoInCommon
					err = db.C("deviceInfos").Find(bson.M{
						"belongTo":   userId,
						"deviceUUID": x.FieldByName("DeviceUUID").String()}).Select(GetSelector(SelectDeviceInfoInCommon)).One(&myDeviceInfo)
					if err == mgo.ErrNotFound {
						// No deviceInfo for this device yet, depending on the situation, let user create on client or get the default back to client.
						err = nil
						myDeviceInfo, err = GetDefaultDeviceInfo(userId, x.FieldByName("DeviceUUID").String(), db)
						if err == mgo.ErrNotFound {
							err = errors.New("No deviceInfo found for this account, please create one on device first.")
						}
					}
					if err == nil {
						// Overwrite client's with server's anyway
						err = SetResBodyPart(v.FieldByName("DeviceInfo"), "DeviceInfo", reflect.ValueOf(myDeviceInfo))
						if err == nil {
							var cardListDB []CardsVerListElement
							cardListDB, err = GetDbCardVerList(db, userId)
							if err == nil {
								fmt.Println("cardListDB length: ", len(cardListDB))
								err = GetCardListDifference(db, cardListDB, x.FieldByName("CardList"), v.FieldByName("CardList"), v.FieldByName("CardToDelete"))
								if err != nil {
									fmt.Println("GetCardListDifference err: ", err.Error())
									_ = ScanToFindAndUpdateEmptyTextId(db, bson.ObjectIdHex(params["user_id"]))
								}
							}
						}
					}
				}
			}
		default:
			err = errors.New("Request not recogniized.")
		}
	} else if m == "DELETE" {
		// One case only: DELETE /users/:user_id/cards/:card_id
		if route == OneCard {
			var t CardInCommon
			var r []CardInCommon
			err = db.C("cards").Find(criteria).One(&t)
			if err == mgo.ErrNotFound {
				// Return ok since no such card exists among non-deleted ones.
				err = nil
			} else if err == nil {
				decisionCode := HandleCardVersionConflict("delete", structFromReq, t)
				if decisionCode == NoConflictOverwriteDB {
					err = db.C("cards").Update(criteria, bson.M{
						"$set": bson.M{"isDeleted": true}})
				} else if decisionCode == ConflictOverwriteClient {
					r = make([]CardInCommon, 0)
					r = append(r, t)
					err = SetResBodyPart(v.FieldByName("Cards"), "Cards", reflect.ValueOf(r))
				}
			}
		} else {
			err = errors.New("Request not recognized.")
		}
	}
	return
}

func SetGetDeviceTokens(userId bson.ObjectId, structFromReq interface{}, db *mgo.Database) (tokens TokensInCommon, err error) {
	x := reflect.ValueOf(structFromReq)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	selector := bson.M{
		"belongTo":   userId,
		"deviceUUID": x.FieldByName("DeviceUUID").String()}
	var t DeviceTokens
	var newId bson.ObjectId
	err = db.C("deviceTokens").Find(selector).One(&t)
	if err == mgo.ErrNotFound {
		err = nil
		// No record found, create new DeviceTokens. SignUp/SignIn
		newId, err = InsertNonDicDB(DeviceTokensNew, structFromReq, db, userId)
		if err == nil {
			err = db.C("deviceTokens").Find(bson.M{"_id": newId}).Select(GetSelector(SelectDeviceTokensInCommon)).One(&tokens)
		}
	} else if err == nil {
		// Record found, update existing record.
		var selector1 bson.M
		selector1, err = UpdateNonDicDB(DeviceTokensUpdateTokens, structFromReq, db, nil, userId)
		if err == nil {
			err = db.C("deviceTokens").Find(selector1).Select(GetSelector(SelectDeviceTokensInCommon)).One(&tokens)
		}
	}
	return
}

func GetDbCardVerList(db *mgo.Database, userId bson.ObjectId) (r []CardsVerListElement, err error) {
	// Shoud identify the err returned for further possibities.
	err = db.C("cards").Find(bson.M{
		"belongTo":  userId,
		"isDeleted": false}).Select(bson.M{
		"_id":       1,
		"versionNo": 1}).All(&r)
	return
}

func GetCardListDifference(db *mgo.Database, cardListDB []CardsVerListElement, cardListReq reflect.Value, cardListRes reflect.Value, cardsToDeleteRes reflect.Value) (err error) {
	// cardListReq: []CardsVerList, cardListRes: []CardInCommon, cardsToDelete: []bson.ObjectId
	// The difference is to overwrite/add docs on client
	// Only non-deleted cards on db here, see PreprocessRequest.
	if !cardListRes.CanSet() || !cardsToDeleteRes.CanSet() {
		return
	}
	cardList := []CardInCommon{}
	cardsToDelete := []bson.ObjectId{}
	fmt.Println("card list db: ", cardListDB)
	for _, x := range cardListDB {
		// Reset the flag
		hasMatch := false
		// Search if the specific id exists
		fmt.Println("id string card db: ", x.Id.Hex())
		fmt.Println("cardListReq.Len: ", cardListReq.Len())
		for y := 0; y < cardListReq.Len(); y++ {
			fmt.Println("id string card req: ", cardListReq.Index(y).FieldByName("Id").String())
			if string(x.Id) == cardListReq.Index(y).FieldByName("Id").String() {
				fmt.Println("Equal")
				hasMatch = true
				if x.VersionNo == cardListReq.Index(y).FieldByName("VersionNo").Int() {
					// Same version, do nothing
				} else {
					// Overwrite the doc on client even if x.VersionNo < cardListReq.Index(y).FieldByName("VersionNo").Int() which indicates something wrong on client.
					cardList, err = PushCardToCardList(db, cardList, x.Id)
					if err != nil {
						return
					}
				}
				break
			}
		}
		// If the specific id not exists, push to list
		if hasMatch == false {
			cardList, err = PushCardToCardList(db, cardList, x.Id)
			fmt.Println("cardList length00000: ", len(cardList))
			if err != nil {
				fmt.Println("cardList length err: ", err.Error())
				return
			}
		}
	}
	fmt.Println("cardList length1111: ", len(cardList))
	if len(cardList) > 0 {
		err = SetResBodyPart(cardListRes, "CardList", reflect.ValueOf(cardList))
	}
	// Just find out which ones are not on server
	for k := 0; k < cardListReq.Len(); k++ {
		noMatch := true
		fmt.Println("cardListReq.Index(k): ", cardListReq.Index(k).FieldByName("Id").String())
		// Search if the specific id exists
		for _, j := range cardListDB {
			fmt.Println("cardListDB: ", string(j.Id))
			if cardListReq.Index(k).FieldByName("Id").String() == string(j.Id) {
				noMatch = false
				break
			}
		}
		if noMatch {
			// If the specific id not exists, push to delete list anyway
			cardsToDelete = append(cardsToDelete, bson.ObjectId(cardListReq.Index(k).FieldByName("Id").String()))
		}
	}
	if len(cardsToDelete) > 0 {
		err = SetResBodyPart(cardsToDeleteRes, "CardToDelete", reflect.ValueOf(cardsToDelete))
	}
	return
}

func PushCardToCardList(db *mgo.Database, cardList []CardInCommon, cardID bson.ObjectId) (newCardList []CardInCommon, err error) {
	var aCard CardInCommon
	err = db.C("cards").Find(bson.M{
		"_id":       cardID,
		"isDeleted": false}).Select(GetSelector(SelectCardInCommon)).One(&aCard)
	fmt.Println("PushCardToCardList card: ", aCard)
	fmt.Println("PushCardToCardList cardId: ", cardID.Hex())
	if err != nil {
		fmt.Println("PushCardToCardList err: ", err.Error())
		return
	}
	if err == nil {
		newCardList = append(cardList, aCard)
	}
	fmt.Println("PushCardToCardList cardList: ", cardList)
	fmt.Println("PushCardToCardList newCardList: ", newCardList)
	// err == mgo.ErrNotFound is still an error here.
	return
}

func SetResBodyPart(partToSet reflect.Value, fieldNameToSet string, valueIn reflect.Value) (err error) {
	if partToSet.CanSet() {
		// v := reflect.ValueOf(valueIn)
		// if v.Kind() == reflect.Ptr {
		// 	v = v.Elem()
		// }
		partToSet.Set(valueIn)
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
		newId, err = InsertNonDicDB(CardNew, structFromReq, db, bson.ObjectIdHex(params["user_id"]))
		if err == nil {
			err = db.C("cards").Find(bson.M{
				"_id": newId}).Select(GetSelector(SelectCardInCommon)).All(&inserted)
		}
	}
	return
}

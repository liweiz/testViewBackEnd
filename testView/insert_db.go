package testView

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

/*
criteria example:
bson.M{
		key: x
		}

cards:
bson.M{
		"belongTo": user.ID,
		"isDeleted": false
		}

deviceInfo:
bson.M{
		"email": user.Email,
		"deviceUUID": targetUUID
		}
*/

const (
	// Each accessToken is valid for 3 hours, (3.6e+12)*3
	// For test purpose, reduce this to 3 mins, (3.6e+12)/20
	tokenExpirationInNanoSec int64 = (3.6e+12) * 3
)

const (
	CardDicNew int = iota
	CardNew
	UserNew
	DeviceTokensNew
	DeviceInfoNew
	RequestProcessedNew
	DicTextContextNew
	DicTextTargetNew
	DicTextTranslationNew
	DicTextDetailNew
)

func InsertNonDicDB(defaultDocType int, structFromReq interface{}, db *mgo.Database, userId bson.ObjectId) (newId bson.ObjectId, err error) {
	var d interface{}
	d, newId, err = PrepareNewNonDicDocDB(defaultDocType, structFromReq, userId)
	// fmt.Println("bson.M to insert 2: ", d)
	if err == nil {
		var name string
		switch defaultDocType {
		case CardNew:
			name = "cards"
		case UserNew:
			name = "users"
		case DeviceTokensNew:
			name = "deviceTokens"
		case DeviceInfoNew:
			name = "deviceInfos"
		case RequestProcessedNew:
			name = "requestsProcessed"
		default:
			err = errors.New("No matched document type for database.")
		}
		if err == nil {
			err = db.C(name).Insert(d)
			if err != nil {
				fmt.Println("bson.M to insert err: ", err.Error())
			}
		}
	}
	return
}

func PrepareNewNonDicDocDB(defaultDocType int, structFromReq interface{}, userId bson.ObjectId) (docToSave interface{}, newId bson.ObjectId, err error) {
	d := reflect.ValueOf(structFromReq)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	newId = bson.NewObjectId()
	switch defaultDocType {
	case CardNew:
		docToSave = bson.M{
			"context":              d.FieldByName("Card").FieldByName("Context").String(),
			"target":               d.FieldByName("Card").FieldByName("Target").String(),
			"translation":          d.FieldByName("Card").FieldByName("Translation").String(),
			"detail":               d.FieldByName("Card").FieldByName("Detail").String(),
			"sourceLang":           d.FieldByName("Card").FieldByName("SourceLang").String(),
			"targetLang":           d.FieldByName("Card").FieldByName("TargetLang").String(),
			"contextDicTextId":     "",
			"detailDicTextId":      "",
			"targetDicTextId":      "",
			"translationDicTextId": "",
			"_id": newId,
			// VersionNo begins from 1, not 0.
			"versionNo":    1,
			"belongTo":     userId,
			"lastModified": time.Now().UnixNano(),
			"collectedAt":  time.Now().UnixNano(),
			"isDeleted":    false}
	case UserNew:
		uuid.SwitchFormat(uuid.Clean, false)
		uniqueUrlCode := uuid.NewV4().String()
		var hashedPassword []byte
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(d.FieldByName("Password").String()), bcrypt.DefaultCost)
		if err == nil {
			docToSave = bson.M{
				// UserInCommon
				"activated":  false,
				"email":      d.FieldByName("Email").String(),
				"_id":        newId,
				"versionNo":  1,
				"isSharing":  false,
				"sourceLang": d.FieldByName("SourceLang").String(),
				"targetLang": d.FieldByName("TargetLang").String(),

				//Non UserInCommon part
				"lastModified":      time.Now().UnixNano(),
				"createdAt":         time.Now().UnixNano(),
				"isDeleted":         false,
				"password":          hashedPassword,
				"activationUrlCode": uniqueUrlCode}
		}

	case DeviceTokensNew:
		accessToken, refreshToken := GenerateTokens(true)
		docToSave = bson.M{
			// TokensInCommon
			"accessToken":         accessToken,
			"refreshToken":        refreshToken,
			"_id":                 newId,
			"belongTo":            userId,
			"deviceUUID":          d.FieldByName("DeviceUUID").String(),
			"accessTokenExpireAt": time.Now().UnixNano() + tokenExpirationInNanoSec,
			"lastModified":        time.Now().UnixNano()}
	// DeviceInfo is created after user is created successfully.
	case DeviceInfoNew:
		docToSave = bson.M{
			// DeviceInfoInCommon
			"_id":        newId,
			"belongTo":   userId,
			"deviceUUID": d.FieldByName("DeviceUUID").String(),
			// These are set by users after a successful signup.
			"sourceLang": d.FieldByName("DeviceInfo").FieldByName("SourceLang").String(),
			"targetLang": d.FieldByName("DeviceInfo").FieldByName("TargetLang").String(),
			"isLoggedIn": true,
			"rememberMe": true,
			"sortOption": d.FieldByName("DeviceInfo").FieldByName("SortOption").String(),

			// Non DeviceInfoInCommon part
			"lastModified": time.Now().UnixNano()}
	case RequestProcessedNew:
		docToSave = bson.M{
			"_id":          newId,
			"belongToUser": userId,
			"deviceUUID":   d.FieldByName("DeviceUUID").String(),
			"requestId":    d.FieldByName("RequestId").String(),
			"timestamp":    time.Now().UnixNano()}
	default:
		err = errors.New("No matched document type for nonDic database.")
	}
	fmt.Println("bson.M to insert: ", docToSave)
	return
}

// case DicTranslation:
// 		resStruct := ResDicResults{}
// 		sourceLang := params["sourcelang"]
// 		targetLang := params["targetlang"]
// 		// words_id from URL here is the text of the words, not id. Coz it's not easy to link user's input to the id.
// 		words := params["words_id"]
// 		// words := bson.ObjectIdHex(params["words_id"])
// 		c := bson.M{
// 			"sourcelang": sourceLang,
// 			"targetlang": targetLang,
// 			// 1: context, 2: target, 3: translation, 4: detail
// 			"textType": 3,
// 			"length":   len(words),
// 			"text":     words,
// 			// Compared with non-deleted cards only to minimize the resource needed.
// 			"isDeleted": false,
// 			"isHidden":  false}
// 		PrepareVehicle(ctx, nil, resStruct, c, "", "")
// 	case DicDetail:
// 		resStruct := ResDicResults{}
// 		sourceLang := params["sourcelang"]
// 		targetLang := params["targetlang"]
// 		translationId := bson.ObjectIdHex(params["translation_id"])
// 		c := bson.M{
// 			"sourcelang": sourceLang,
// 			"targetlang": targetLang,
// 			"textType":   4,
// 			"belongTo":   translationId,
// 			// Compared with non-deleted cards only to minimize the resource needed.
// 			"isDeleted": false,
// 			"isHidden":  false}
// 		PrepareVehicle(ctx, nil, resStruct, c, "", "")
// 	case DicContext:
// 		resStruct := ResDicResults{}
// 		sourceLang := params["sourcelang"]
// 		targetLang := params["targetlang"]
// 		contextId := bson.ObjectIdHex(params["context_id"])
// 		c := bson.M{
// 			"sourcelang": sourceLang,
// 			"targetlang": targetLang,
// 			"textType":   1,
// 			"belongTo":   contextId,
// 			// Compared with non-deleted cards only to minimize the resource needed.
// 			"isDeleted": false,
// 			"isHidden":  false}
// 		PrepareVehicle(ctx, nil, resStruct, c, "", "")

// Dic in db can only be created. The only update action is to update its LastModified field, which is done through another func.
func InsertDicDB(defaultDocType int, cardInCommon interface{}, db *mgo.Database, parentId bson.ObjectId) (err error) {
	var d interface{}
	var now int64
	d, now, err = PrepareDicDocForDB(defaultDocType, cardInCommon, parentId)
	switch defaultDocType {
	case DicTextContextNew:
		if err == nil {
			err = db.C("dicTextContexts").Insert(d)
			if err == nil {
				// Update Detail level parent
				var grandParentId bson.ObjectId
				grandParentId, err = UpdateDicParentLastModified(db, "dicTextDetails", parentId, now)
				if err == nil {
					var greatGrandParentId bson.ObjectId
					greatGrandParentId, err = UpdateDicParentLastModified(db, "dicTextTranslations", grandParentId, now)
					if err == nil {
						_, err = UpdateDicParentLastModified(db, "dicTextTargets", greatGrandParentId, now)
					}
				}
			}
		}
	case DicTextTargetNew:
		if err == nil {
			err = db.C("dicTextTargets").Insert(d)
		}
	case DicTextTranslationNew:
		if err == nil {
			err = db.C("dicTextTranslations").Insert(d)
			if err == nil {
				_, err = UpdateDicParentLastModified(db, "dicTextTargets", parentId, now)
			}
		}
	case DicTextDetailNew:
		if err == nil {
			err = db.C("dicTextDetails").Insert(d)
			if err == nil {
				var grandParentId bson.ObjectId
				grandParentId, err = UpdateDicParentLastModified(db, "dicTextTranslations", parentId, now)
				if err == nil {
					_, err = UpdateDicParentLastModified(db, "dicTextTargets", grandParentId, now)
				}
			}
		}
	default:
		err = errors.New("No matched document type for database.")
	}
	return
}

func PrepareDicDocForDB(defaultDocType int, card interface{}, parentId bson.ObjectId) (docToSave interface{}, now int64, err error) {
	// This is executed after the card is created/updated. So we get the info from the card in db instead of the request body.
	// structFromReq here is used as the card struct in db.
	d := reflect.ValueOf(card)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	// 1: context, 2: target, 3: translation, 4: detail
	var textType int
	var textTypeName string
	switch defaultDocType {
	case DicTextContextNew:
		textType = 1
		textTypeName = "Context"
	case DicTextTargetNew:
		textType = 2
		textTypeName = "Target"
	case DicTextTranslationNew:
		textType = 3
		textTypeName = "Translation"
	case DicTextDetailNew:
		textType = 4
		textTypeName = "Detail"
	default:
		err = errors.New("No matched document type for Dic database.")
	}
	if err != nil {
		newId := bson.NewObjectId()
		now = time.Now().UnixNano()
		docToSave = bson.M{
			"_id":            newId,
			"sourceLangCode": AssignCodeToLang(d.FieldByName("SourceLang").String()),
			"targetLangCode": AssignCodeToLang(d.FieldByName("TargetLang").String()),
			"sourceLang":     d.FieldByName("SourceLang").String(),
			"targetLang":     d.FieldByName("TargetLang").String(),
			// 1: context, 2: target, 3: translation, 4: detail
			"textType":                     textType,
			"text":                         d.FieldByName(textTypeName).String(),
			"textLength":                   len(d.FieldByName(textTypeName).String()),
			"belongToParent":               parentId,
			"noOfUsersPickedThisUpFromDic": 0,
			"isDeleted":                    false,
			"isHidden":                     false,
			"createdAt":                    now,
			"lastModified":                 now,
			"createdBy":                    d.FieldByName("BelongTo").String(),
			"childrenLastUpdatedAt":        now}
	}
	return
}

func UpdateDicParentLastModified(db *mgo.Database, collectionName string, parentId bson.ObjectId, now int64) (grandParentId bson.ObjectId, err error) {
	var r struct {
		BelongToParent bson.ObjectId `bson:"belongToParent" json:"belongToParent"`
		LastModified   int64         `bson:"lastModified" json:"lastModified"`
	}
	err = db.C(collectionName).Find(bson.M{
		"_id": parentId}).Select(bson.M{
		"belongToParent": 1,
		"lastModified":   1}).One(&r)
	if err == nil {
		if now > r.LastModified {
			err = db.C(collectionName).Update(bson.M{
				"_id": parentId}, bson.M{
				"$set": bson.M{
					"lastModified": now}})
			if err == nil {
				grandParentId = r.BelongToParent
			}
		}
	}
	return
}

func FillResourceInStruct(structToFill interface{}, structToCheck interface{}) {
	x := reflect.ValueOf(structToFill)
	if x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	iMax := x.NumField()
	y := reflect.ValueOf(structToCheck)
	if y.Kind() == reflect.Ptr {
		y = y.Elem()
	}
	z := x.Type()
	for i := 0; i < iMax; i++ {
		zz := z.Field(i).Name
		if y.FieldByName(zz) != reflect.ValueOf(nil) {
			x.Field(i).Set(y.FieldByName(zz))
		}
	}
}

func PrepareDicResultDocForDB(defaultDocType int, structFromReq interface{}, params martini.Params) (docToSave interface{}, newId bson.ObjectId, err error) {
	d := reflect.ValueOf(structFromReq)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	switch defaultDocType {

	case DeviceInfoUpdateDicTier2:
		docToSave = bson.M{
			"$set": bson.M{
				"dicTier2": ""}}
	default:
		err = errors.New("No matched document type for dicResult database.")
	}
	return
}

// func GetSortedDicResult(db *mgo.Database, dicName string, idToSearch bson.ObjectId, sortOption string) (result []DicTextInRes) {
// 	if sortOption
// 	db.C(dicName).Find(bson.M{"belongToParent": idToSearch}).Sort().Select(bson.M{"_id": 1, "text": 1, "childrenLastUpdatedAt": 1}).All(&result)
// }

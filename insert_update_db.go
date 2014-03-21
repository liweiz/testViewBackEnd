package testView

import (
        "string"
        "reflect"
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "github.com/codegangsta/martini"
        "errors"
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
	CardDicNew 							= iota
	CardNew
	CardUpdate
	UserNew
	UserUpdateActivation
	UserUpdateEmail
	UserUpdatePassword
	DeviceTokensNew
	DeviceTokensUpdateTokens
	DeviceInfoNew
	DeviceInfoUpdateSortOption
	DeviceInfoUpdateLang
	RequestProcessedNew
	DicTextContextNew
	DicTextTargetNew
	DicTextTranslationNew
	DicTextDetailNew
)

func InsertUpdateNonDicDB(defaultDocType int, structFromReq interface{}, db *mgo.Database, params *martini.Params) (err error) {
	var d interface{}
	d, err = PrepareNonDicDocForDB(defaultDocType, structFromReq, params)
	if defaultDocType == UserUpdateActivation || defaultDocType == UserUpdateEmail || defaultDocType == UserUpdatePassword {
		UserUpdate := UserUpdateActivation
		defaultDocType = UserUpdate
	} else if defaultDocType == DeviceInfoUpdateSortOption || defaultDocType == DeviceInfoUpdateLang {
		DeviceInfoUpdate := DeviceInfoUpdateSortOption
		defaultDocType = DeviceInfoUpdate
	}
	switch defaultDocType {
	case CardNew:
		if err == nil {
			err = db.C("cards").Insert(d)
		}
	case CardUpdate:
		if err == nil {
			err = db.C("cards").Update(bson.M{
				"_id": bson.ObjectIdHex(params[card_id]),
				"belongTo": bson.ObjectIdHex(params[user_id])
				}, d)
		}
	case UserNew:
		if err == nil {
			err = db.C("users").Insert(d)
		}
	case UserUpdate:
		if err == nil {
			err = db.C("users").Update(bson.M{
				"_id": bson.ObjectIdHex(params[user_id])
				}, d)
		}
	case DeviceTokensNew:
		if err == nil {
			err = db.C("deviceTokens").Insert(d)
		}
	case DeviceTokensUpdateTokens:
		if err == nil {
			x := reflect.ValueOf(structFromReq)
			if x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			err = db.C("deviceTokens").Update(bson.M{
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"belongTo": bson.ObjectIdHex(params[user_id])
				}, d)
		}
	case DeviceInfoNew:
		if err == nil {
			err = db.C("deviceInfos").Insert(d)
		}
	case DeviceInfoUpdate:
		if err == nil {
			x := reflect.ValueOf(structFromReq)
			if x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			err = db.C("deviceInfos").Update(bson.M{
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"belongTo": bson.ObjectIdHex(params[user_id])
				}, d)
		}
	case RequestProcessedNew:
		if err == nil {
			err = db.C("requestsProcessed").Insert(d)
		}
	default:
		err = errors.New("No matched document type for database.")
	}
	return
}

func PrepareNonDicDocForDB(defaultDocType int, structFromReq interface{}, params *martini.Params) (docToSave interface{}, err error) {
	d := reflect.ValueOf(structFromReq)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	switch defaultDocType {
	case CardNew:
		newId = bson.NewObjectId()
		docToSave := bson.M{
			"context": d.FieldByName("Card").FieldByName("Context").String(),
			"target": d.FieldByName("Card").FieldByName("Target").String(),
			"translation": d.FieldByName("Card").FieldByName("Translation").String(),
			"detail": d.FieldByName("Card").FieldByName("Detail").String(),
			"sourceLang": d.FieldByName("Card").FieldByName("SourceLang").String(),
			"targetLang": d.FieldByName("Card").FieldByName("TargetLang").String(),
			"_id": bson.NewObjectId(),
			// VersionNo begins from 1, not 0.
			"versionNo": 1,
			"belongTo": bson.ObjectIdHex(params[user_id]),
			"lastModified": time.Now().UnixNano(),
			"collectedAt": time.Now().UnixNano(),
			"isDeleted": false
		}
	case CardUpdate:
		docToSave := bson.M{
			"$set": bson.M{
				"context": d.FieldByName("Card").FieldByName("Context").String(),
				"target": d.FieldByName("Card").FieldByName("Target").String(),
				"translation": d.FieldByName("Card").FieldByName("Translation").String(),
				"detail": d.FieldByName("Card").FieldByName("Detail").String(),
				"lastModified": time.Now().UnixNano()
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case UserNew:
		newId = bson.NewObjectId()
		docToSave := bson.M{
			// UserInCommon
			"activated": false,
			"email": d.FieldByName("Email").String(),
			"_id": bson.NewObjectId(),
			"versionNo": 1,

			//Non UserInCommon part
			"lastModified": time.Now().UnixNano(),
			"createdAt": time.Now().UnixNano(),
			"isDeleted": false,
			"password": d.FieldByName("Password").String()
		}
	case UserUpdateActivation:
		docToSave := bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"activated": true
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case UserUpdateEmail:
		docToSave := bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"email": d.FieldByName("NewEmail").String()
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case UserUpdatePassword:
		docToSave := bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"password": d.FieldByName("NewPassword").String()
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case DeviceTokensNew:
		newId = bson.NewObjectId()
		accessToken, refreshToken := GenerateTokens(true)
		docToSave := bson.M{
			// TokensInCommon
			"accessToken": accessToken,
			"refreshToken": refreshToken,

			"_id": newId,
			"belongTo": userId,
			"deviceUUID": d.FieldByName("DeviceUUID").String(),
			// Each accessToken is valid for 3 hours
			"accessTokenExpireAt": time.Now().UnixNano() + (3.6e+12) * 3,
			"lastModified": time.Now().UnixNano()
		}
	case DeviceTokensUpdateTokens:
		accessToken, refreshToken := GenerateTokens(true)
		docToSave := bson.M{
			"$set": bson.M{
				"accessToken": accessToken,
				"refreshToken": refreshToken,
				"accessTokenExpireAt": time.Now().UnixNano() + (3.6e+12) * 3,
				"lastModified": time.Now().UnixNano()
			}
		}
	// DeviceInfo is created after user is created successfully.
	case DeviceInfoNew:
		newId = bson.NewObjectId()
		docToSave := bson.M{
			// DeviceInfoInCommon
			"_id": newId,
			"belongTo": bson.ObjectIdHex(params[user_id]),
			"deviceUUID": d.FieldByName("DeviceUUID").String(),
			// These are set by users after a successful signup.
			"sourceLang": "",
			"targetLang": "",
			"isLoggedIn": true,
			"rememberMe": true,
			"sortOption": "ByLastModifiedDescending",
			"versionNo": 1,

			// Non DeviceInfoInCommon part
			"lastModified": time.Now().UnixNano()
		}
	case DeviceInfoUpdateSortOption:
		docToSave := bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"sortOption": d.FieldByName("DeviceInfo").FieldByName("SortOption").String()
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case DeviceInfoUpdateLang:
		docToSave := bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"sourceLang": d.FieldByName("DeviceInfo").FieldByName("SourceLang").String(),
				"targetLang": d.FieldByName("DeviceInfo").FieldByName("TargetLang").String()
			}
			"$inc": bson.M{
				"versionNo": 1
			}
		}
	case RequestProcessedNew:
		newId = bson.NewObjectId()
		docToSave := bson.M{
			"_id": newId,
			"belongToUser": bson.ObjectIdHex(params[user_id]),
			"deviceUUID": d.FieldByName("DeviceUUID").String(),
			"requestVersionNo": d.FieldByName("RequestVersionNo").Int(),
			"timestamp": time.Now().UnixNano()
		}
	default:
		err = errors.New("No matched document type for nonDic database.")
	}
	return
}

// Dic in db can only be created. The only update action is to update its LastModified field, which is done through another func.
func InsertDicDB(defaultDocType int, cardInCommon interface{}, db *mgo.Database, parentId bson.ObjectId) (err error) {
	var d interface{}
	var now int64
	var idInserted bson.ObjectId
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
				var grandParentId bson.ObjectId
				grandParentId, err = UpdateDicParentLastModified(db, "dicTextTargets", parentId, now)
			}
		}
	case DicTextDetailNew:
		if err == nil {
			err = db.C("dicTextDetails").Insert(d)
			if err == nil {
				var grandParentId bson.ObjectId
				grandParentId, err = UpdateDicParentLastModified(db, "dicTextTranslations", parentId, now)
				if err == nil {
					var greatGrandParentId bson.ObjectId
					greatGrandParentId, err = UpdateDicParentLastModified(db, "dicTextTargets", grandParentId, now)
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
		docToSave := bson.M{
			"_id": newId,
			"sourceLang": d.FieldByName("SourceLang").String(),
			"targetLang": d.FieldByName("TargetLang").String(),
			// 1: context, 2: target, 3: translation, 4: detail
			"textType": textType,
			"text": d.FieldByName(textTypeName).String(),
			"textLength": len(d.FieldByName(textTypeName).String()),
			"belongToParent": parentId,
			"isDeleted": false,
			"isHidden": false,
			"createdAt": now,
			"lastModified": now,
			"createdBy": d.FieldByName("BelongTo").String(),
			"childrenLastUpdatedAt": now
		}
	}
	return
}

func UpdateDicParentLastModified(db *mgo.Database, collectionName string, parentId bson.ObjectId, now int64) (grandParentId bson.ObjectId, err error) {
	var r struct {
		BelongToParent bson.ObjectId 'bson:"belongToParent" json:"belongToParent"'
		LastModified int64 'bson:"lastModified" json:"lastModified"'
	}
	err = db.C(collectionName).Find(bson.M{
		"_id": parentId
		}).Select(bson.M{
			"belongToParent": 1,
			"lastModified": 1
			}).One(&r)
	if err == nil {
		if now > r.LastModified {
			err = db.C(collectionName).Update(bson.M{
			"_id": parentId
			}, bson.M{
				"$set": bson.M{
					"lastModified": now,
				})
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
    	x = x..Elem()
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

// Insert with uniqueness check, user and id can be nil in some cases
// doc represents different things in different actions.
func CRUDOneDoc(action string, db *mgo.Database, doc interface{}, criteria bson.M, params *martini.Params) (err error) {
	// Create and Update are the most frequent operations, so put them in the front
	// Check uniqueness first since there is no need to create/update a non-unique doc. Get doc's type name at the same time
	
	v := reflect.ValueOf(doc)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var docTypeName string
	if t.Name() == "Card" {
		docTypeName = "cards"
	} else if t.Name() == "User" {
		docTypeName = "users"
	} else if t.Name() == "CardInDic" {
		docTypeName = "cardInDics"
	} else {
		err = errors.New("Unidentified document to CRUD")
		return
	}
	
	if action == "create" {
		isUnique := CheckDocUniqueness(db, doc, params)
		if isUnique {
			err = db.C(docTypeName).Insert(doc)
		} else {

		}
	} else if action == "create" && !isUnique {
		err = errors.New("")
	}
	if action == "update" {
		err = db.C(docTypeName).Insert(doc)
	}
	// There is no need to check uniqueness for deletion and reading
	if action == "read" {
		err = db.C(docTypeName).Find(criteria).One(doc)
	}
	if action =="flagDelete" {
		// Not real deletion, just change the isDeleted flag. The isDeleted flag as to be false to proceed.
		err = db.C(docTypeName).Update(criteria, bson.M{
			"isDeleted": true
			})
	}
	return
}
package testView

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

const (
	CardUpdate int = iota
	UserUpdateActivation
	UserUpdateEmail
	UserAddPasswordUrlCode
	UserRemovePasswordUrlCodeRoutinely
	UserUpdatePassword
	DeviceTokensUpdateTokens
	DeviceInfoUpdateSortOption
	DeviceInfoUpdateLang
	DeviceInfoUpdateDicTier2
	DeviceInfoUpdateDicTier3
	DeviceInfoUpdateDicTier4
)

func UpdateNonDicDB(defaultDocType int, structFromReq interface{}, db *mgo.Database, params martini.Params, userId bson.ObjectId) (selector bson.M, err error) {
	var d interface{}
	d, err = PrepareUpdateNonDicDocDB(defaultDocType, structFromReq)
	var UserUpdate int
	var DeviceInfoUpdate int
	if defaultDocType == UserUpdateActivation || defaultDocType == UserUpdateEmail || defaultDocType == UserUpdatePassword || defaultDocType == UserAddPasswordUrlCode || defaultDocType == UserRemovePasswordUrlCodeRoutinely {
		UserUpdate = UserUpdateActivation
		defaultDocType = UserUpdate
	} else if defaultDocType == DeviceInfoUpdateSortOption || defaultDocType == DeviceInfoUpdateLang {
		DeviceInfoUpdate = DeviceInfoUpdateSortOption
		defaultDocType = DeviceInfoUpdate
	}
	if err == nil {
		var name string
		switch defaultDocType {
		case CardUpdate:
			selector = bson.M{
				"_id":      bson.ObjectIdHex(params["card_id"]),
				"belongTo": userId}
			name = "cards"
		case UserUpdate:
			selector = bson.M{"_id": userId}
			name = "users"
		case DeviceTokensUpdateTokens:
			x := reflect.ValueOf(structFromReq)
			if x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			selector = bson.M{
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"belongTo":   userId}
			name = "deviceTokens"
		case DeviceInfoUpdate:
			x := reflect.ValueOf(structFromReq)
			if x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			selector = bson.M{
				"deviceUUID": x.FieldByName("DeviceUUID").String(),
				"belongTo":   userId}
			name = "deviceInfos"
		default:
			err = errors.New("No matched document type for database.")
		}
		if err == nil {
			err = db.C(name).Update(selector, d)
		}
	}
	return
}

// Only updating on XxxInCommon triggers versionNo increment. Coz versionNo is only useful to compare common content on both server and client.
func PrepareUpdateNonDicDocDB(defaultDocType int, structFromReq interface{}) (docToSave interface{}, err error) {
	d := reflect.ValueOf(structFromReq)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	switch defaultDocType {
	case CardUpdate:
		docToSave = bson.M{
			"$set": bson.M{
				"context":      d.FieldByName("Card").FieldByName("Context").String(),
				"target":       d.FieldByName("Card").FieldByName("Target").String(),
				"translation":  d.FieldByName("Card").FieldByName("Translation").String(),
				"detail":       d.FieldByName("Card").FieldByName("Detail").String(),
				"lastModified": time.Now().UnixNano()},
			"$inc": bson.M{
				"versionNo": 1}}
	case UserUpdateActivation:
		docToSave = bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"activated":    true},
			"$inc": bson.M{
				"versionNo": 1}}
	case UserUpdateEmail:
		docToSave = bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"email":        d.FieldByName("NewEmail").String()},
			"$inc": bson.M{
				"versionNo": 1}}
	case UserAddPasswordUrlCode:
		uuid.SwitchFormat(uuid.Clean, false)
		uniqueUrlCode := uuid.NewV4().String()
		docToSave = bson.M{
			"$push": bson.M{
				"passwordResettingUrlCodes": bson.M{"passwordResettingUrlCode": uniqueUrlCode,
					"timeStamp": time.Now().UnixNano()}},
			"$set": bson.M{
				"lastModified": time.Now().UnixNano()}}
	case UserRemovePasswordUrlCodeRoutinely:
		// A uniqueUrl is only valid for one hour.
		y := time.Now().UnixNano() - (3.6e+12)
		docToSave = bson.M{
			"$pull": bson.M{
				"passwordResettingUrlCodes": bson.M{"timeStamp": bson.M{"$lt": y}}},
			"$set": bson.M{
				"lastModified": time.Now().UnixNano()}}
	case UserUpdatePassword:
		var hashedPassword []byte
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(d.FieldByName("NewPassword").String()), bcrypt.DefaultCost)
		if err == nil {
			docToSave = bson.M{
				"$set": bson.M{
					"lastModified": time.Now().UnixNano(),
					"password":     hashedPassword}}
		}
	case DeviceTokensUpdateTokens:
		accessToken, refreshToken := GenerateTokens(true)
		docToSave = bson.M{
			"$set": bson.M{
				"accessToken":         accessToken,
				"refreshToken":        refreshToken,
				"accessTokenExpireAt": time.Now().UnixNano() + (3.6e+12)*3,
				"lastModified":        time.Now().UnixNano()}}
	case DeviceInfoUpdateSortOption:
		docToSave = bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"sortOption":   d.FieldByName("DeviceInfo").FieldByName("SortOption").String()}}
	case DeviceInfoUpdateLang:
		docToSave = bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"sourceLang":   d.FieldByName("DeviceInfo").FieldByName("SourceLang").String(),
				"targetLang":   d.FieldByName("DeviceInfo").FieldByName("TargetLang").String()}}
	default:
		err = errors.New("No matched document type for nonDic database.")
	}
	return
}

// This only happens when urlCode is not expired.
func RemoveUserPasswordUrlCodeOne(db *mgo.Database, params martini.Params, userId bson.ObjectId) (err error) {
	err = db.C("users").Update(bson.M{"_id": userId}, bson.M{
		"$pull": bson.M{
			"passwordResettingUrlCodes": bson.M{"passwordResettingUrlCode": params["password_resetting_code"]}},
		"$set": bson.M{
			"lastModified": time.Now().UnixNano()}})
	return
}

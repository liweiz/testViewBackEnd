package testView

import (
	"errors"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

const (
	CardUpdate int = iota
	UserUpdateActivation
	UserUpdateEmail
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
	if defaultDocType == UserUpdateActivation || defaultDocType == UserUpdateEmail || defaultDocType == UserUpdatePassword {
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
	case UserUpdatePassword:
		docToSave = bson.M{
			"$set": bson.M{
				"lastModified": time.Now().UnixNano(),
				"password":     d.FieldByName("NewPassword").String()},
			"$inc": bson.M{
				"versionNo": 1}}
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

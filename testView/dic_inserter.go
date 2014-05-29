package testView

import (
	// "github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

func ScanToFindAndUpdateEmptyTextId(db *mgo.Database, userId bson.ObjectId) (err error) {
	var a []Card
	err = db.C("cards").Find(bson.M{"belongTo": userId}).All(&a)
	if err == nil && len(a) > 0 {
		for _, x := range a {
			sLangCode := AssignCodeToLang(x.SourceLang)
			tLangCode := AssignCodeToLang(x.TargetLang)
			var targetId bson.ObjectId
			var translationId bson.ObjectId
			var detailId bson.ObjectId
			var contextId bson.ObjectId
			if x.TargetDicTextId == "" {
				var d1 DicTextInCommon
				d1, err = GetSearchedDicTextObj(db, x.Target, 2, sLangCode, tLangCode)
				if err == mgo.ErrNotFound {
					targetId = bson.NewObjectId()
					err = db.C("dicTexts").Insert(bson.M{
						"_id":                   targetId,
						"sourceLangCode":        sLangCode,
						"targetLangCode":        tLangCode,
						"sourceLang":            x.SourceLang,
						"targetLang":            x.TargetLang,
						"textType":              2,
						"text":                  x.Target,
						"textLength":            len(x.Target),
						"belongToParent":        "",
						"noOfUsersHavingThis":   1,
						"isDeleted":             false,
						"isHidden":              false,
						"createdAt":             time.Now().UnixNano(),
						"lastModified":          time.Now().UnixNano(),
						"createdBy":             userId,
						"childrenLastUpdatedAt": time.Now().UnixNano(),
					})
					if err != nil {
						return
					}
				} else if err == nil {
					targetId = d1.Id
				} else {
					return
				}
				err = db.C("cards").Update(bson.M{"_id": x.Id}, bson.M{"$set": bson.M{"targetDicTextId": targetId}})
				if err != nil {
					return
				}
			} else {
				targetId = x.TargetDicTextId
			}
			if x.TranslationDicTextId == "" {
				var d2 DicTextInCommon
				d2, err = GetSearchedDicTextObj(db, x.Translation, 3, sLangCode, tLangCode)
				if err == mgo.ErrNotFound {
					translationId = bson.NewObjectId()
					err = db.C("dicTexts").Insert(bson.M{
						"_id":                   translationId,
						"sourceLangCode":        sLangCode,
						"targetLangCode":        tLangCode,
						"sourceLang":            x.SourceLang,
						"targetLang":            x.TargetLang,
						"textType":              3,
						"text":                  x.Translation,
						"textLength":            len(x.Translation),
						"belongToParent":        targetId,
						"noOfUsersHavingThis":   1,
						"isDeleted":             false,
						"isHidden":              false,
						"createdAt":             time.Now().UnixNano(),
						"lastModified":          time.Now().UnixNano(),
						"createdBy":             userId,
						"childrenLastUpdatedAt": time.Now().UnixNano(),
					})
					if err != nil {
						return
					} else {
						err = db.C("dicTexts").Update(bson.M{"_id": targetId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
						if err != nil {
							return
						}
					}
				} else if err == nil {
					translationId = d2.Id
				} else {
					return
				}
				err = db.C("cards").Update(bson.M{"_id": x.Id}, bson.M{"$set": bson.M{"translationDicTextId": translationId}})
				if err != nil {
					return
				}
			} else {
				translationId = x.TranslationDicTextId
			}
			if x.DetailDicTextId == "" {
				var d3 DicTextInCommon
				d3, err = GetSearchedDicTextObj(db, x.Detail, 4, sLangCode, tLangCode)
				if err == mgo.ErrNotFound {
					detailId = bson.NewObjectId()
					err = db.C("dicTexts").Insert(bson.M{
						"_id":                   detailId,
						"sourceLangCode":        sLangCode,
						"targetLangCode":        tLangCode,
						"sourceLang":            x.SourceLang,
						"targetLang":            x.TargetLang,
						"textType":              4,
						"text":                  x.Detail,
						"textLength":            len(x.Detail),
						"belongToParent":        translationId,
						"noOfUsersHavingThis":   1,
						"isDeleted":             false,
						"isHidden":              false,
						"createdAt":             time.Now().UnixNano(),
						"lastModified":          time.Now().UnixNano(),
						"createdBy":             userId,
						"childrenLastUpdatedAt": time.Now().UnixNano(),
					})
					if err != nil {
						return
					} else {
						err = db.C("dicTexts").Update(bson.M{"_id": translationId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
						if err != nil {
							return
						} else {
							err = db.C("dicTexts").Update(bson.M{"_id": targetId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
							if err != nil {
								return
							}
						}
					}
				} else if err == nil {
					detailId = d3.Id
				} else {
					return
				}
				err = db.C("cards").Update(bson.M{"_id": x.Id}, bson.M{"$set": bson.M{"detailDicTextId": detailId}})
				if err != nil {
					return
				}
			} else {
				detailId = x.DetailDicTextId
			}
			if x.ContextDicTextId == "" {
				var d4 DicTextInCommon
				d4, err = GetSearchedDicTextObj(db, x.Context, 1, sLangCode, tLangCode)
				if err == mgo.ErrNotFound {
					contextId = bson.NewObjectId()
					err = db.C("dicTexts").Insert(bson.M{
						"_id":                   contextId,
						"sourceLangCode":        sLangCode,
						"targetLangCode":        tLangCode,
						"sourceLang":            x.SourceLang,
						"targetLang":            x.TargetLang,
						"textType":              1,
						"text":                  x.Context,
						"textLength":            len(x.Context),
						"belongToParent":        detailId,
						"noOfUsersHavingThis":   1,
						"isDeleted":             false,
						"isHidden":              false,
						"createdAt":             time.Now().UnixNano(),
						"lastModified":          time.Now().UnixNano(),
						"createdBy":             userId,
						"childrenLastUpdatedAt": time.Now().UnixNano(),
					})
					if err != nil {
						return
					} else {
						err = db.C("dicTexts").Update(bson.M{"_id": detailId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
						if err != nil {
							return
						} else {
							err = db.C("dicTexts").Update(bson.M{"_id": translationId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
							if err != nil {
								return
							} else {
								err = db.C("dicTexts").Update(bson.M{"_id": targetId}, bson.M{"$set": bson.M{"childrenLastUpdatedAt": time.Now().UnixNano()}})
								if err != nil {
									return
								}
							}
						}
					}
				} else if err == nil {
					contextId = d4.Id
				} else {
					return
				}
				err = db.C("cards").Update(bson.M{"_id": x.Id}, bson.M{"$set": bson.M{"contextDicTextId": contextId}})
				if err != nil {
					return
				}
			} else {
				contextId = x.ContextDicTextId
			}
		}
	}
	return
}

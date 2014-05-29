package testView

import (
	"encoding/json"
	"errors"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AssignCodeToLang(langName string) (code int) {
	switch langName {
	case "العربية":
		code = 1
	case "Български":
		code = 2
	case "Català":
		code = 3
	case "Česky":
		code = 4
	case "Deutsch":
		code = 5
	case "Eesti":
		code = 6
	case "English":
		code = 7
	case "Español":
		code = 8
	case "Euskara":
		code = 9
	case "فارسی":
		code = 10
	case "Français":
		code = 11
	case "Galego":
		code = 12
	case "한국어":
		code = 13
	case "Հայերեն":
		code = 14
	case "Bahasa Indonesia":
		code = 15
	case "Italiano":
		code = 16
	case "עברית":
		code = 17
	case "Latviešu":
		code = 18
	case "Magyar":
		code = 19
	case "മലയാളം":
		code = 20
	case "Bahasa Melayu":
		code = 21
	case "Nederlands":
		code = 22
	case "日本語":
		code = 23
	case "Norsk bokmål":
		code = 24
	case "Polski":
		code = 25
	case "Português":
		code = 26
	case "Română":
		code = 27
	case "Русский":
		code = 28
	case "Српски / srpski":
		code = 29
	case "Suomi":
		code = 30
	case "Svenska":
		code = 31
	case "தமிழ்":
		code = 32
	case "ไทย":
		code = 33
	case "Türkçe":
		code = 34
	case "Українська":
		code = 35
	case "Tiếng Việt":
		code = 36
	case "简体中文":
		code = 37
	case "繁體中文":
		code = 38
	default:
		// not found
		code = 0
	}
	return
}

func DicTextSearcher() martini.Handler {
	return func(db *mgo.Database, req *http.Request, params martini.Params, logger *log.Logger, rw http.ResponseWriter) {
		var s string
		p := PathInDic{}
		p.SourceLangCode, _ = strconv.Atoi(params["source_lang_code"])
		p.TargetLangCode, _ = strconv.Atoi(params["target_lang_code"])
		reqStruct := &ReqDicText{}
		resStruct := &ResDicResultsText{}
		err := GetStructFromReq(req, reqStruct)
		if err == nil {
			if len(reqStruct.WordsText) > 0 {
				var d DicTextInCommon
				d, err = GetSearchedDicTextObj(db, reqStruct.WordsText, 2, p.SourceLangCode, p.TargetLangCode)
				if err == nil {
					p.TargetDicTextId = d.Id
					var r []DicTextInRes
					r, err = p.GetRangeSearchResultInDicById(db, 3, "", reqStruct.SortOption, reqStruct.IsAscending)
					if err == nil {
						resStruct.TextType = 3
						resStruct.TopLevelTextId = p.TargetDicTextId
						resStruct.Results = r
						rw.Header().Set("Content-Type", "application/json")
						var j []byte
						j, err = json.Marshal(resStruct)
						if err == nil {
							// Response size, any usage???
							_, err = rw.Write(j)
						}
						if err != nil {
							s = strings.Join([]string{"Failed to generate response, but request has been successfully processed by server.", err.Error()}, "=> ")
						}
					}
				} else if err == mgo.ErrNotFound {
					err = errors.New("No such text found in dic.")
				}
			} else {
				err = errors.New("Non empty text needed for searching.")
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

func DicIdSearcher() martini.Handler {
	return func(db *mgo.Database, req *http.Request, params martini.Params, logger *log.Logger, rw http.ResponseWriter) {
		var s string
		p := PathInDic{}
		p.SourceLangCode, _ = strconv.Atoi(params["source_lang_code"])
		p.TargetLangCode, _ = strconv.Atoi(params["target_lang_code"])
		reqStruct := &ReqDicId{}
		resStruct := &ResDicResultsId{}
		err := GetStructFromReq(req, reqStruct)
		if err == nil {
			var parentD DicText
			err = db.C("dicTexts").Find(bson.M{"_id": reqStruct.ParentId}).One(&parentD)
			if err == nil {
				var level int
				switch parentD.TextType {
				case 2:
					level = 3
				case 3:
					level = 4
				case 4:
					level = 1
				default:
					err = errors.New("Incorrect parent textType.")
				}
				if err == nil {
					var r []DicTextInRes
					r, err = p.GetRangeSearchResultInDicById(db, level, reqStruct.LastId, reqStruct.SortOption, reqStruct.IsAscending)
					if err == nil {
						resStruct.TextType = level
						resStruct.Results = r
						rw.Header().Set("Content-Type", "application/json")
						var j []byte
						j, err = json.Marshal(resStruct)
						if err == nil {
							// Response size, any usage???
							_, err = rw.Write(j)
						}
						if err != nil {
							s = strings.Join([]string{"Failed to generate response, but request has been successfully processed by server.", err.Error()}, "=> ")
						}
					}
				}
			} else if err == mgo.ErrNotFound {
				err = errors.New("No such id found in dic with this textType.")
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

func GetSearchedDicTextObj(db *mgo.Database, t string, textType int, sCode int, tCode int) (d DicTextInCommon, err error) {
	// Get id for searched text
	err = db.C("dicTexts").Find(bson.M{
		"sourceLangCode": sCode,
		"targetLangCode": tCode,
		// 1: context, 2: target, 3: translation, 4: detail
		"textType":   textType,
		"text":       t,
		"textLength": len(t),
		"isDeleted":  false,
		"isHidden":   false,
	}).Select(GetSelector(SelectDicInCommon)).One(&d)
	return
}

type PathInDic struct {
	SourceLangCode       int           `bson:"sourceLangCode" json:"sourceLangCode"`
	TargetLangCode       int           `bson:"targetLangCode" json:"targetLangCode"`
	SourceLang           string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang           string        `bson:"targetLang" json:"targetLang"`
	Context              string        `bson:"context" json:"context"`
	ContextDicTextId     bson.ObjectId `bson:"contextDicTextId" json:"contextDicTextId"`
	Detail               string        `bson:"detail" json:"detail"`
	DetailDicTextId      bson.ObjectId `bson:"detailDicTextId" json:"detailDicTextId"`
	Target               string        `bson:"target" json:"target"`
	TargetDicTextId      bson.ObjectId `bson:"targetDicTextId" json:"targetDicTextId"`
	Translation          string        `bson:"translation" json:"translation"`
	TranslationDicTextId bson.ObjectId `bson:"translationDicTextId" json:"translationDicTextId"`
}

func (p *PathInDic) SetLangInPathInDic(source string, target string) {
	p.SourceLang = source
	p.TargetLang = target
	p.SourceLangCode = AssignCodeToLang(source)
	p.TargetLangCode = AssignCodeToLang(target)
}

// Sort option: time created/time modified/most collected/alphabet
func (p *PathInDic) GetRangeSearchResultInDicById(db *mgo.Database, typeCode int, lastId bson.ObjectId, sortBy string, isAscending bool) (r []DicTextInRes, err error) {
	var parentId bson.ObjectId
	switch typeCode {
	// 1: context, 2: target, 3: translation, 4: detail
	case 1:
		parentId = p.DetailDicTextId
	case 2:
		parentId = ""
	case 3:
		parentId = p.TargetDicTextId
	case 4:
		parentId = p.TranslationDicTextId
	}
	var noToSkip int
	if lastId == "" {
		noToSkip = 0
	} else {
		var aDic DicText
		err = db.C("dicTexts").FindId(lastId).One(&aDic)
		if err == nil {
			noToSkip, err = p.CountLastIdPosition(db, parentId, typeCode, sortBy, GetConditionForTarget(sortBy, isAscending, aDic))
			if err != nil {
				return
			}
		}
	}
	err = db.C("dicTexts").Find(bson.M{
		"sourceLangCode": p.SourceLangCode,
		"targetLangCode": p.TargetLangCode,
		"belongToParent": parentId,
		"textType":       typeCode,
		"isDeleted":      false,
		"isHidden":       false,
	}).Sort(GetSortCondition(sortBy, isAscending)).Select(bson.M{
		"_id":      1,
		"text":     1,
		"textType": 1}).Skip(noToSkip).Limit(30).All(&r)
	return
}

func (p *PathInDic) CountLastIdPosition(db *mgo.Database, parentId bson.ObjectId, typeCode int, sortBy string, condition bson.M) (result int, err error) {
	var d []DicText
	err = db.C("dicTexts").Find(bson.M{
		"sourceLangCode": p.SourceLangCode,
		"targetLangCode": p.TargetLangCode,
		"belongToParent": parentId,
		"textType":       typeCode,
		"isDeleted":      false,
		"isHidden":       false,
		sortBy:           condition,
	}).All(&d)
	result = len(d)
	return
}

func GetSortCondition(sortBy string, isAscending bool) (condition string) {
	switch sortBy {
	case "createdAt":
		if isAscending {
			condition = "createdAt"
		} else {
			condition = "-createdAt"
		}
	case "childrenLastUpdatedAt":
		if isAscending {
			condition = "childrenLastUpdatedAt"
		} else {
			condition = "-childrenLastUpdatedAt"
		}
	case "noOfUsersHavingThis":
		if isAscending {
			condition = "noOfUsersHavingThis"
		} else {
			condition = "-noOfUsersHavingThis"
		}
	}
	return
}

func GetConditionForTarget(sortBy string, isAscending bool, doc DicText) (condition bson.M) {
	switch sortBy {
	case "createdAt":
		if isAscending {
			condition = bson.M{"$lte": doc.CreatedAt}
		} else {
			condition = bson.M{"$gte": doc.CreatedAt}
		}
	case "childrenLastUpdatedAt":
		if isAscending {
			condition = bson.M{"$lte": doc.ChildrenLastUpdatedAt}
		} else {
			condition = bson.M{"$gte": doc.ChildrenLastUpdatedAt}
		}
	case "noOfUsersPickedThisUpFromDic":
		if isAscending {
			condition = bson.M{"$lte": doc.NoOfUsersHavingThis}
		} else {
			condition = bson.M{"$gte": doc.NoOfUsersHavingThis}
		}
	}
	return
}

// IT IS BETTER TO LET CLIENT REMEMBER THE LAST ID INSTEAD. NEED TO BE MODIFIED LATER.
// func SetDeviceSearchResult(db *mgo.Database, deviceInfoId bson.ObjectId, r []DicTextInRes, tierName string) error {
// 	return db.C("deviceInfos").Update(bson.M{"_id": deviceInfoId}, bson.M{"$set": bson.M{tierName: r}})
// }

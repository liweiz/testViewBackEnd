package testView

import (
	"labix.org/v2/mgo/bson"
)

type CardToDeleteInRes struct {
	ID bson.ObjectId `json:"id"`
	// VersionNo int64         `json:"versionNo"`
}

type ResSignUpOrIn struct {
	User   UserInCommon   `json:"user"`
	Tokens TokensInCommon `json:"tokens"`

	// Add lang pair here due to no deviceInfo for now
	SourceLang string `bson:"sourceLang" json:"sourceLang"`
	TargetLang string `bson:"targetLang" json:"targetLang"`
}

// RequestVersionNo is used to solve the online/offline sync issue. Tokens is an exception since there is no need to sync tokens, only getting new tokens is needed. Therefore, VersionNo is not necessary here, either.
// Request related to Tokens is proceeded indenpendently with user/card.
type ResTokensOnly struct {
	/*
	   Two possible situations:
	   1. user info already exists on client
	   2. no user info on client
	   There is no easy way to tell which situation is on client. So leave that part for sync request sent by client. A sync request follows the successful login response immediately.
	*/
	UserId string         `json:"userId"`
	Tokens TokensInCommon `json:"tokens"`
	// No RequestVersionNo here. Tokens are not compared by VersionNo/RequestVersionNo. They are compared directly.
	// No need to have user_id as well. Coz a given set of tokens are only valid with one device and user. If anything goes wrong, simply ask the user to login will do.
}

// Sync user
type ResUser struct {
	User UserInCommon `json:"user"`
}

type ResUserNewEmail struct {
	NewEmail string `json:"newEmail"`
}

type ResResetPassword struct {
	PasswordIsResetted bool `json:"passwordIsResetted"`
}

type ResActivation struct {
	IsActivatedThisTime bool `json:"isActivatedThisTime"`
}

type ResDeviceInfo struct {
	DeviceInfo DeviceInfoInCommon `json:"deviceInfo"`
}

// For full card info returned in res
type ResCards struct {
	Cards []CardInCommon `json:"cards"`
}

type ResSync struct {
	User         UserInCommon       `json:"user"`
	DeviceInfo   DeviceInfoInCommon `json:"deviceInfo"`
	CardList     []CardInCommon     `json:"cardList"`
	CardToDelete []bson.ObjectId    `json:"cardToDelete"`
}

type DicTextInRes struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	Text     string        `bson:"text" json:"text"`
	TextType int           `bson:"textType" json:"textType"`
}

type ResDicResultsText struct {
	TextType       int            `json:"textType"`
	TopLevelTextId bson.ObjectId  `json:"topLevelTextId"`
	Results        []DicTextInRes `json:"results"`
}

type ResDicResultsId struct {
	// Level: translation/detail/context
	TextType int            `json:"textType"`
	Results  []DicTextInRes `json:"results"`
}

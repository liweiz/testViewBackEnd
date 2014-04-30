package testView

import (
	"labix.org/v2/mgo/bson"
)

// Only non-GET requests come with RequestId since nothing is changed with a GET request.

// No use, coz they can be stored in Authorization header.
type ReqSignUpOrIn struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceUUID string `json:"deviceUUID"`
}

type ReqForgotPassword struct {
	Email string `json:"email"`
}

type ReqResetPassword struct {
	NewPassword string `json:"newPassword"`
}

type ReqUpdateEmail struct {
	NewEmail string `json:"newEmail"`
}

type ReqRenewTokens struct {
	// UserId is get from params in url.
	DeviceUUID string         `json:"deviceUUID"`
	Tokens     TokensInCommon `json:"tokens"`
}

type ReqUser struct {
	User UserInCommon `json:"user"`
}

// RequestId appears together with client UUID to locate the list to put RequestId into.
type ReqDeviceInfo struct {
	RequestId  string                `json:"requestId"`
	DeviceUUID string                `json:"deviceUUID"`
	DeviceInfo DeviceInfoInCommonNew `json:"deviceInfo"`
}

// Request for card(s) does not provide user id since it can be got from param in URL.
// Request for updating card does not include card's id since it can be got from param in URL.
// So the only difference between new and update is the request URL.
type ReqCard struct {
	RequestId     string          `json:"requestId"`
	DeviceUUID    string          `json:"deviceUUID"`
	Card          CardInCommonNew `json:"card"`
	CardVersionNo int64           `json:"cardVersionNo"`
}

type CardsVerListElement struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	VersionNo int64         `bson:"versionNo" json:"versionNo"`
}

type ReqSync struct {
	DeviceUUID string `json:"deviceUUID"`
	// A map with id and versionNo
	// Deleted cards are not included in the map. If it is not deleted on server, just add it back to the client again witht the sync response.
	CardList []CardsVerListElement `json:"cardList"`
}

type ReqDicText struct {
	WordsText   string `json:"wordsText"`
	SortOption  string `json:"sortOption"`
	IsAscending bool   `json:"isAscending"`
}

type ReqDicId struct {
	ParentId    bson.ObjectId `json:"parentId"`
	LastId      bson.ObjectId `json:"lastId"`
	SortOption  string        `json:"sortOption"`
	IsAscending bool          `json:"isAscending"`
}

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

// RequestId appears together with client UUID to locate the list to put RequestId into.
type ReqUser struct {
	RequestId  string       `json:"requestId"`
	DeviceUUID string       `json:"deviceUUID"`
	User       UserInCommon `json:"user"`
}

type ReqDeviceInfoNew struct {
	RequestId  string                `json:"requestId"`
	DeviceUUID string                `json:"deviceUUID"`
	DeviceInfo DeviceInfoInCommonNew `json:"deviceInfo"`
}

type ReqDeviceInfo struct {
	RequestId  string             `json:"requestId"`
	DeviceUUID string             `json:"deviceUUID"`
	DeviceInfo DeviceInfoInCommon `json:"deviceInfo"`
}

// Request for card(s) do not provide user id since it can be got from param in URL
type ReqCardNew struct {
	RequestId  string          `json:"requestId"`
	DeviceUUID string          `json:"deviceUUID"`
	Card       CardInCommonNew `json:"card"`
}

type ReqCard struct {
	RequestId  string       `json:"requestId"`
	DeviceUUID string       `json:"deviceUUID"`
	Card       CardInCommon `json:"card"`
}

type ReqCards struct {
	RequestId  string         `json:"requestId"`
	DeviceUUID string         `json:"deviceUUID"`
	Cards      []CardInCommon `json:"cards"`
}

type CardsVerList struct {
	Id        bson.ObjectId `json:"id"`
	VersionNo int64         `json:"versionNo"`
}

type ReqSync struct {
	RequestId  string             `json:"requestId"`
	DeviceUUID string             `json:"deviceUUID"`
	User       UserInCommon       `json:"user"`
	DeviceInfo DeviceInfoInCommon `json:"deviceInfo"`
	// A map with id and versionNo
	// Deleted cards are not included in the map. If it is not deleted on server, just add it back to the client again witht the sync response.
	CardList []CardsVerList `json:"cardList"`
}

type ReqWords struct {
	WordsText string `json:"wordsText"`
}

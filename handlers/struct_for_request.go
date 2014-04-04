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

type ReqResetPassword struct {
	NewPassword string
}

type ReqUpdateEmail struct {
	NewEmail string
}

type ReqRenewTokens struct {
	// UserId is get from params in url.
	DeviceUUID string         `json:"deviceUUID"`
	Tokens     TokensInCommon `json:"tokens"`
}

// RequestId appears together with client UUID to locate the list to put RequestId into.
type ReqUser struct {
	RequestId  string
	DeviceUUID string
	User       UserInCommon
}

type ReqDeviceInfo struct {
	RequestId  string
	DeviceUUID string
	DeviceInfo DeviceInfoInCommon
}

// Request for card(s) do not provide user id since it can be got from param in URL
type ReqCard struct {
	RequestId  string
	DeviceUUID string
	Card       CardInCommon
}

type ReqCards struct {
	RequestId  string
	DeviceUUID string
	Cards      []CardInCommon
}

type CardsVerList struct {
	Id        bson.ObjectId
	VersionNo int64
}

type ReqSync struct {
	RequestId  string
	DeviceUUID string
	User       UserInCommon
	DeviceInfo DeviceInfoInCommon
	// A map with id and versionNo
	// Deleted cards are not included in the map. If it is not deleted on server, just add it back to the client again witht the sync response.
	CardList []CardsVerList
}

type ReqWords struct {
	WordsText string
}

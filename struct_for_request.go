package testView

// Only non-GET requests come with ReqVerNo since nothing is changed with a GET request.

// No use, coz they can be stored in Authorization header.
type ReqSignUpOrIn struct {
	Email string
	Password string
	DeviceUUID string
}

type ReqResetPassword struct {
	NewPassword string
}

type ReqUpdateEmail struct {
	NewEmail string
}

type ReqRenewTokens struct {
	DeviceUUID string
	Tokens TokensInCommon
}

// ReqVerNo appears together with client UUID to locate the list to put ReqVerNo into.
type ReqUser struct {
	ReqVerNo int64
	DeviceUUID string
	User UserInCommon
}

type ReqDeviceInfo struct {
	ReqVerNo int64
	DeviceUUID string
	DeviceInfo DeviceInfoInCommon
}

// Request for card(s) do not provide user id since it can be got from param in URL
type ReqCard struct {
	ReqVerNo int64
	DeviceUUID string
	Card CardInCommon
}

type ReqCards struct {
	ReqVerNo int64
	DeviceUUID string
	Cards []CardInCommon
}

type CardsVerList struct {
	Id bson.ObjectId
	VersionNo int64
}

type ReqSync struct {
	ReqVerNo int64
	DeviceUUID string
	User UserInCommon
	DeviceInfo DeviceInfoInCommon
	// A map with id and versionNo
	// Deleted cards are not included in the map. If it is not deleted on server, just add it back to the client again witht the sync response.
	CardList []CardsVerList
}

type ReqWords struct {
	WordsText string
}
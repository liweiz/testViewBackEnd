package testView

import (
	"labix.org/v2/mgo/bson"
)

// gorp will pass time.Time fields through to the database/sql driver, but note that the behavior of this type varies across database drivers.

// MySQL users should be especially cautious. See: https://github.com/ziutek/mymysql/pull/77

// To avoid any potential issues with timezone/DST, consider using an integer field for time data and storing UNIX time.

// Time is stored as int64 in UnixNano, which returns a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC..

const (
	SelectUserInCommon int = iota
	SelectDicInCommon
	SelectDeviceTokensInCommon
	SelectDeviceInfoInCommon
	SelectCardInCommon
)

func GetSelector(option int) (r bson.M) {
	switch option {
	case SelectUserInCommon:
		r = bson.M{
			"activated": 1,
			"email":     1,
			"_id":       1,
			"versionNo": 1}
	case SelectDicInCommon:
		r = bson.M{
			"_id":                          1,
			"textType":                     1,
			"text":                         1,
			"textLength":                   1,
			"noOfUsersPickedThisUpFromDic": 1,
			"createdAt":                    1,
			"childrenLastUpdatedAt":        1}
	case SelectDeviceTokensInCommon:
		r = bson.M{
			"accessToken":  1,
			"refreshToken": 1}
	case SelectDeviceInfoInCommon:
		r = bson.M{
			"dicTier2":     0,
			"dicTier3":     0,
			"dicTier4":     0,
			"lastModified": 0}
	case SelectCardInCommon:
		r = bson.M{
			"isDeleted": 0}
	}
	return
}

type UserInCommon struct {
	Activated  bool          `bson:"activated" json:"activated"`
	Email      string        `bson:"email" json:"email"`
	Id         bson.ObjectId `bson:"_id" json:"_id"`
	VersionNo  int64         `bson:"versionNo" json:"versionNo"`
	SourceLang string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang string        `bson:"targetLang" json:"targetLang"`
}

type PasswordResettingUrlCodePair struct {
	PasswordResettingUrlCode string `bson:"passwordResettingUrlCode" json:"passwordResettingUrlCode"`
	TimeStamp                int64  `bson:"timeStamp" json:"timeStamp"`
}

type User struct {
	Activated bool          `bson:"activated" json:"activated"`
	Email     string        `bson:"email" json:"email"`
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	VersionNo int64         `bson:"versionNo" json:"versionNo"`

	// Add lang pair to User instead of deviceInfo to restrict each user can only have one lang pair. So we removed the multi lang pair function for one user to reduce the chance that the incorrect lang pair set to card. By abandoning deviceInfo, we also removed the sort option sync and other device specific settings. User always gets default sort option when first login occurs, and manual change has to be done to change the settings. If one user wants to have different lang pair, another account will be needed in this case.
	SourceLang string `bson:"sourceLang" json:"sourceLang"`
	TargetLang string `bson:"targetLang" json:"targetLang"`

	// Server side only
	LastModified      int64  `bson:"lastModified" json:"lastModified"`
	CreatedAt         int64  `bson:"createdAt" json:"createdAt"`
	Password          []byte `bson:"password" json:"password"`
	IsDeleted         bool   `bson:"isDeleted" json:"isDeleted"`
	ActivationUrlCode string `bson:"activationUrlCode" json:"activationUrlCode"`
	// PasswordResettingUrls is a url/timestamp pair. Expired urls will be evaluated and cleaned up each time a new password-resetting request is received by server.
	PasswordResettingUrlCodes []*PasswordResettingUrlCodePair `bson:"passwordResettingUrlCodes" json:"passwordResettingUrlCodes"`
	// For versionNo calculation, newVersionNo = HighestVersionNo + 1
	// RequestProcessed and HasCards have no effect on VersionNo. In other words, no VersionNo change when either of these two changes.
}

type TokensInCommon struct {
	// When being used, the _id is obtained through device_id in url. E.g., /user/:user_id/deviceinfos/device_id
	AccessToken  string `bson:"accessToken" json:"accessToken"`
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
}

type DeviceTokens struct {
	AccessToken  string `bson:"accessToken" json:"accessToken"`
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`

	Id                  bson.ObjectId `bson:"_id" json:"_id"`
	BelongTo            bson.ObjectId `bson:"belongTo" json:"belongTo"`
	DeviceUUID          string        `bson:"deviceUUID" json:"deviceUUID"`
	AccessTokenExpireAt int64         `bson:"accessTokenExpireAt" json:"accessTokenExpireAt"`
	LastModified        int64         `bson:"lastModified" json:"lastModified"`
}

/*
Different cases:
1. Sign up:
	Client gets user in response first and then provides interface for language pair setting. Use default sortOption if user does not change it specifically.
2. Sign in first time on a given device. No previous local file:
	Client
*/

type DeviceInfoInCommonNew struct {
	DeviceUUID string `bson:"deviceUUID" json:"deviceUUID"`
	// LanguagePair selection is done when user signs up. User will be asked to set the pair after a successful signup but before they can use the app. The initially selected pair is stored and used as default. The change operation will be provided in future releases. Not a problem currently.
	// LanguagePair selection is a device-specific option. Each device can has its own preference
	SourceLang string `bson:"sourceLang" json:"sourceLang"`
	TargetLang string `bson:"targetLang" json:"targetLang"`
	// SortOption/IsLoggedIn/RememberMe are device-specific as well.
	SortOption string `bson:"sortOption" json:"sortOption"`
	IsLoggedIn bool   `bson:"isLoggedIn" json:"isLoggedIn"`
	RememberMe bool   `bson:"rememberMe" json:"rememberMe"`
}

type DeviceInfoInCommon struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`

	BelongTo   bson.ObjectId `bson:"belongTo" json:"belongTo"`
	DeviceUUID string        `bson:"deviceUUID" json:"deviceUUID"`
	SourceLang string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang string        `bson:"targetLang" json:"targetLang"`
	SortOption string        `bson:"sortOption" json:"sortOption"`
	IsLoggedIn bool          `bson:"isLoggedIn" json:"isLoggedIn"`
	RememberMe bool          `bson:"rememberMe" json:"rememberMe"`
}

// Server side only
type DeviceInfo struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`

	BelongTo   bson.ObjectId `bson:"belongTo" json:"belongTo"`
	DeviceUUID string        `bson:"deviceUUID" json:"deviceUUID"`
	SourceLang string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang string        `bson:"targetLang" json:"targetLang"`
	SortOption string        `bson:"sortOption" json:"sortOption"`
	IsLoggedIn bool          `bson:"isLoggedIn" json:"isLoggedIn"`
	RememberMe bool          `bson:"rememberMe" json:"rememberMe"`

	LastModified int64 `bson:"lastModified" json:"lastModified"`
	// Save search result's last object's objectId for each tier and serve by pagination. No tier1 here since tier1 is just one words. Overwrite corresponding tier when new search is triggered by client. Meanwhile, clear up the sibling tiers if there is any. The result is sorted by the device`s sortOption stored.
	// Change on these does not update lastModified in DeviceInfo.
	// DicTier2 bson.ObjectId `bson:"dicTier2" json:"dicTier2"`
	// DicTier3 bson.ObjectId `bson:"dicTier3" json:"dicTier3"`
	// DicTier4 bson.ObjectId `bson:"dicTier4" json:"dicTier4"`
}

// Server side only
/*
There should be a requestVersionNo list on each client as well.
Each client has its own requestVerisonList.
Sync flow:
1: A non-sync request with reqVerNo(not all req is sent with this) sent by a client.
2: The client adds its requestVersionNo to the list. The list has requestVersionNo and isDone.
2.1: The client does not receive response/receives non-200 response from server. Send the same request again till 200 response received.
3: The server receives the request. Find out if it has been successfully processed with its RequestProcessed. If not, process the request. Otherwise, send 200 response back to the client, no body.
3.1: After the successful process of the request. Send 200 response with the results in body. No need to add the request to RequestProcessed.
3.2: When process is not successful, send the err message back to the client.
2.2: A 200 response received by client. Set the requestVersionNo off the list. If there is body, modify the local data with the results from server. If no body, which is only to let the client know the server has successfully processed the request, do nothing and let the sync request update the corresponding local data.
4: When requestVersionNo pool on client is cleared, send sync request. Before receiving the response of the sync request, no other request is sent out by the client. This is to avoid potential data inconsistence.
4.1: No/non-200 response received, sync is not successful. If there are other requests awaiting, send them first, back into the loop till the request pool is cleared again then send the updated sync request and the old version of sync request is dropped from the list since it is no longer needed. Otherwise, resend sync request till 200 response is received.
4.2: 200 response is received. Update local data according to the body.

After the corresponding response received by the client, it change the isDone to true. Any sync request is triggered after all existing requestVersionNos has their isDone set to be true. Forcing to sync leads to checking the list and resending the requests that isDone is set as false.
Server keeps the response info in RequestProcessed to handle duplicated requests (requests with same requestVersionNo). After the sync request is received, the server clear the previous response info since they are not needed any more. The way to trigger sync request on client indicates sync request sent only when all previous requests` responses have been received. It`s actually a circle. Between every two adjacent sync requests, unDone requests are sent and the later sync request is triggered indicates previous requests are all been processed. The earlier one indicates the unDone requests start to be accumulated. The possible operation before the sync response is received by client makes the client drop the response and trigger another sync request after the new operation is done. The dropped sync request needs to be set as isDone = true on both server and client.
*/
// A requestProcessed is created after the response is successfully sent. For now, only needed for card since other records are device and account specific. In other words, each device and account combination has its own data store on server. The User model on server always overwrites clients'.
type RequestProcessed struct {
	Id           bson.ObjectId `bson:"_id" json:"_id"`
	BelongToUser bson.ObjectId `bson:"belongToUser" json:"belongToUser"`
	// Use UUID instead of BelongToDevice since it`s easier to locate a RequestProcessed with UUID. Save the step to use UUID to get the deviceId.
	DeviceUUID string `bson:"deviceUUID" json:"deviceUUID"`
	// Use randomly generated string like token as RequestId instead of int64 to avoid potential duplicated RequestId. This field was originally named RequestVersionNo which takes an int64.
	RequestId string `bson:"requestId" json:"requestId"`
	// IsDone bool `bson:"isDone" json:"isDone"`
	Timestamp int64 `bson:"timestamp" json:"timestamp"`
}

type CardInCommonNew struct {
	Context     string `bson:"context" json:"context"`
	Detail      string `bson:"detail" json:"detail"`
	SourceLang  string `bson:"sourceLang" json:"sourceLang"`
	Target      string `bson:"target" json:"target"`
	TargetLang  string `bson:"targetLang" json:"targetLang"`
	Translation string `bson:"translation" json:"translation"`
}

type CardInCommon struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`

	Context      string        `bson:"context" json:"context"`
	Detail       string        `bson:"detail" json:"detail"`
	SourceLang   string        `bson:"sourceLang" json:"sourceLang"`
	Target       string        `bson:"target" json:"target"`
	TargetLang   string        `bson:"targetLang" json:"targetLang"`
	Translation  string        `bson:"translation" json:"translation"`
	VersionNo    int64         `bson:"versionNo" json:"versionNo"`
	CollectedAt  int64         `bson:"collectedAt" json:"collectedAt"`
	LastModified int64         `bson:"lastModified" json:"lastModified"`
	BelongTo     bson.ObjectId `bson:"belongTo" json:"belongTo"`
}

type Card struct {
	Id                   bson.ObjectId `bson:"_id" json:"_id"`
	SourceLang           string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang           string        `bson:"targetLang" json:"targetLang"`
	Context              string        `bson:"context" json:"context"`
	Detail               string        `bson:"detail" json:"detail"`
	Target               string        `bson:"target" json:"target"`
	Translation          string        `bson:"translation" json:"translation"`
	ContextDicTextId     bson.ObjectId `bson:"contextDicTextId" json:"contextDicTextId"`
	DetailDicTextId      bson.ObjectId `bson:"detailDicTextId" json:"detailDicTextId"`
	TargetDicTextId      bson.ObjectId `bson:"targetDicTextId" json:"targetDicTextId"`
	TranslationDicTextId bson.ObjectId `bson:"translationDicTextId" json:"translationDicTextId"`
	VersionNo            int64         `bson:"versionNo" json:"versionNo"`
	CollectedAt          int64         `bson:"collectedAt" json:"collectedAt"`
	LastModified         int64         `bson:"lastModified" json:"lastModified"`
	BelongTo             bson.ObjectId `bson:"belongTo" json:"belongTo"`

	// Server side only
	IsDeleted bool `bson:"isDeleted" json:"isDeleted"`
}

type DicTextInCommon struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`
	// 1: context, 2: target, 3: translation, 4: detail
	TextType              int    `bson:"textType" json:"textType"`
	Text                  string `bson:"text" json:"text"`
	TextLength            int    `bson:"textLength" json:"textLength"`
	NoOfUsersHavingThis   int64  `bson:"noOfUsersHavingThis" json:"noOfUsersHavingThis"`
	CreatedAt             int64  `bson:"createdAt" json:"createdAt"`
	ChildrenLastUpdatedAt int64  `bson:"childrenLastUpdatedAt" json:"childrenLastUpdatedAt"`
}

type DicText struct {
	Id             bson.ObjectId `bson:"_id" json:"_id"`
	SourceLangCode int           `bson:"sourceLangCode" json:"sourceLangCode"`
	TargetLangCode int           `bson:"targetLangCode" json:"targetLangCode"`
	SourceLang     string        `bson:"sourceLang" json:"sourceLang"`
	TargetLang     string        `bson:"targetLang" json:"targetLang"`
	// 1: context, 2: target, 3: translation, 4: detail
	TextType              int           `bson:"textType" json:"textType"`
	Text                  string        `bson:"text" json:"text"`
	TextLength            int           `bson:"textLength" json:"textLength"`
	BelongToParent        bson.ObjectId `bson:"belongToParent" json:"belongToParent"`
	NoOfUsersHavingThis   int64         `bson:"noOfUsersHavingThis" json:"noOfUsersHavingThis"`
	IsDeleted             bool          `bson:"isDeleted" json:"isDeleted"`
	IsHidden              bool          `bson:"isHidden" json:"isHidden"`
	CreatedAt             int64         `bson:"createdAt" json:"createdAt"`
	LastModified          int64         `bson:"lastModified" json:"lastModified"`
	CreatedBy             bson.ObjectId `bson:"createdBy" json:"createdBy"`
	ChildrenLastUpdatedAt int64         `bson:"childrenLastUpdatedAt" json:"childrenLastUpdatedAt"`
}

/*
Structure for Dic in db:
												   DicTextTarget
											//						\\
							DicTextTranslation							DicTextTranslation
							//				\\
				DicTextDetail				DicTextDetail
				//			\\
	DicTextContext			DicTextContext
*/

type DicTextContext struct {
	DicText
}

type DicTextTarget struct {
	DicText
}

type DicTextTranslation struct {
	DicText
}

type DicTextDetail struct {
	DicText
}

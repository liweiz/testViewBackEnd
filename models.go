package testView

import {
	"time"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
}

// gorp will pass time.Time fields through to the database/sql driver, but note that the behavior of this type varies across database drivers.

// MySQL users should be especially cautious. See: https://github.com/ziutek/mymysql/pull/77

// To avoid any potential issues with timezone/DST, consider using an integer field for time data and storing UNIX time.

// Time is stored as int64 in UnixNano, which returns a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC..

/*

*/

SelectUserInCommon := bson.M{
	"activated": 1,
	"email": 1,
	"_id": 1,
	"versionNo": 1}

type UserInCommon struct {
	Activated bool 'bson:"activated" json:"activated"'
	Email string 'bson:"email" json:"email"'	
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	VersionNo int64 'bson:"versionNo" json:"versionNo"'
}

type User struct {
	UserInCommon

	// Server side only
	LastModified int64 'bson:"lastModified" json:"lastModified"'
	CreatedAt int64 'bson:"createdAt" json:"createdAt"'
	Password string 'bson:"password" json:"password"'
	IsDeleted bool 'bson:"isDeleted" json:"isDeleted"'
	// For versionNo calculation, newVersionNo = HighestVersionNo + 1
	// RequestProcessed and HasCards have no effect on VersionNo. In other words, no VersionNo change when either of these two changes.
}

SelectDeviceTokensInCommon := bson.M{
	"accessToken": 1,
	"refreshToken": 1}

type TokensInCommon struct {
	// When being used, the _id is obtained through device_id in url. E.g., /user/:user_id/deviceinfos/device_id
	AccessToken string 'bson:"accessToken" json:"accessToken"'
	RefreshToken string 'bson:"refreshToken" json:"refreshToken"'
}

type DeviceTokens struct {
	TokensInCommon
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	BelongTo bson.ObjectId 'bson:"belongTo" json:"belongTo"'
	DeviceUUID string 'bson:"deviceUUID" json:"deviceUUID"'
	AccessTokenExpireAt int64 'bson:"accessTokenExpireAt" json:"accessTokenExpireAt"'
	LastModified int64 'bson:"lastModified" json:"lastModified"'
}

/*
Different cases:
1. Sign up:
	Client gets user in response first and then provides interface for language pair setting. Use default sortOption if user does not change it specifically.
2. Sign in first time on a given device. No previous local file:
	Client
*/
SelectDeviceInfoInCommon := bson.M{
	"lastModified": 0}

type DeviceInfoInCommon struct {
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	BelongTo bson.ObjectId 'bson:"belongTo" json:"belongTo"'
	DeviceUUID string 'bson:"deviceUUID" json:"deviceUUID"'
	// LanguagePair selection is done when user signs up. User will be asked to set the pair after a successful signup but before they can use the app. The initially selected pair is stored and used as default. The change operation will be provided in future releases. Not a problem currently.
	// LanguagePair selection is a device-specific option. Each device can has its own preference
	SourceLang string 'bson:"sourceLang" json:"sourceLang"'
	TargetLang string 'bson:"targetLang" json:"targetLang"'
	// SortOption/IsLoggedIn/RememberMe are device-specific as well.
	SortOption string 'bson:"sortOption" json:"sortOption"'
	IsLoggedIn bool 'bson:"isLoggedIn" json:"isLoggedIn"'
	RememberMe bool 'bson:"rememberMe" json:"rememberMe"'
}

// Server side only
type DeviceInfo struct {
	DeviceInfoInCommon
	// Updating tokens does not change VersionNo here.
	LastModified int64 'bson:"lastModified" json:"lastModified"'
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
3.1: After the successful process of the request. Send 200 response with the results in body. Add the request to RequestProcessed.
3.2: When process is not successful, send the err message back tot eh client.
2.2: A 200 response received by client. Set the requestVersionNo off the list. If there is body, modify the local data with the results from server. If no body, do nothing and let the sync request update the corresponding local data.
4: When requestVersionNo pool on client is cleared, send sync request. Before receiving the response of the sync request, no other request is sent out by the client. This is to avoid potential data inconsistence.
4.1: No/non-200 response received, sync is not successful. If there are other requests awaiting, send them first, back into the loop till the request pool is cleared again then send the updated sync request and the old version of sync request is dropped from the list since it is no longer needed. Otherwise, resend sync request till 200 response is received.
4.2: 200 response is received. Update local data according to the body.

After the corresponding response received by the client, it change the isDone to true. Any sync request is triggered after all existing requestVersionNos has their isDone set to be true. Forcing to sync leads to checking the list and resending the requests that isDone is set as false.
Server keeps the response info in RequestProcessed to handle duplicated requests (requests with same requestVersionNo). After the sync request is received, the server clear the previous response info since they are not needed any more. The way to trigger sync request on client indicates sync request sent only when all previous requests' responses have been received. It's actually a circle. Between every two adjacent sync requests, unDone requests are sent and the later sync request is triggered indicates the all set of previous requests. The earlier one indicates the unDone requests start to be accumulated. The possible operation before the sync response is received by client makes the client drop the response and trigger another sync request after the new operation is done. The dropped sync request needs to be set as isDone = true on both server and client.
*/
// A requestProcessed is created after the response is successfully sent.
type RequestProcessed struct {
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	BelongToUser bson.ObjectId 'bson:"belongToUser" json:"belongToUser"'
	// Use UUID instead of BelongToDevice since it's easier to locate a RequestProcessed with UUID. Save the step to use UUID to get the deviceId.
	DeviceUUID string 'bson:"deviceUUID" json:"deviceUUID"'
	RequestVersionNo int64 'bson:"requestVersionNo" json:"requestVersionNo"'
	// IsDone bool 'bson:"isDone" json:"isDone"'
	Timestamp int64 'bson:"timestamp" json:"timestamp"'
}

SelectCardInCommon := bson.M{
	"belongTo": 0,
	"isDeleted": 0}

type CardInCommon {
	Context string 'bson:"context" json:"context"'
	Detail string 'bson:"detail" json:"detail"'
	SourceLang string 'bson:"sourceLang" json:"sourceLang"'
	Target string 'bson:"target" json:"target"'
	TargetLang string 'bson:"targetLang" json:"targetLang"'
	Translation string 'bson:"translation" json:"translation"'
	// HasTag []string
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	VersionNo int64 'bson:"versionNo" json:"versionNo"'
	CollectedAt int64 'bson:"collectedAt" json:"collectedAt"'
	LastModified int64 'bson:"lastModified" json:"lastModified"'
}

type Card struct {
	CardInCommon

	// Server side only
	BelongTo bson.ObjectId 'bson:"belongTo" json:"belongTo"'
	IsDeleted bool 'bson:"isDeleted" json:"isDeleted"'
}

/*

*/
type DicText struct {
	Id bson.ObjectId 'bson:"_id" json:"_id"'
	SourceLang string 'bson:"sourceLang" json:"sourceLang"'
	TargetLang string 'bson:"targetLang" json:"targetLang"'
	// 1: context, 2: target, 3: translation, 4: detail
	TextType int 'bson:"textType" json:"textType"'
	Text string 'bson:"text" json:"text"'
	TextLength int 'bson:"textLength" json:"textLength"'
	BelongToParent bson.ObjectId 'bson:"belongToParent" json:"belongToParent"'

	IsDeleted bool 'bson:"isDeleted" json:"isDeleted"'
	IsHidden bool 'bson:"isHidden" json:"isHidden"'
	CreatedAt int64 'bson:"createdAt" json:"createdAt"'
	LastModified int64 'bson:"lastModified" json:"lastModified"'
	CreatedBy bson.ObjectId 'bson:"createdBy" json:"createdBy"'
	ChildrenLastUpdatedAt int64 'bson:"childrenLastUpdatedAt" json:"childrenLastUpdatedAt"'
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
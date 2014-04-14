package testView

import (
	"labix.org/v2/mgo/bson"
)

/*
Response body structure(JSON):

RESTful for single card:
Request:
{
    "requestVersionNo": number,

    "user": {
    	"email": string,
    	"activated": boolean,
    	"isLoggedIn": boolean,
    	"rememberMe": boolean,
    	"sortOption": string,
    	"LastModified": number,
    	"id": string,
    	"versionNo": number
    },

    "card": {
		"belongTo": string,
    	"collectedAt": number,
    	"context": string,
    	"createdAt": number,
    	"createdBy": string,
        "detail": string,
        "sourceLang": string,
        "target": string,
        "targetLang": string,
        "translation": string,
        "HasTag": array,
        "collectedBy": string,
        "lastModified": number,
        "id": string,
        "versionNo": number,
    },
}
New: send full card
Update: send full card
Delete: send id and versionNo of card only

Response:

Sync request:
{
    "requestVersionNo": number,

    "user": {
        "email": string,
        "activated": boolean,
        "isLoggedIn": boolean,
        "rememberMe": boolean,
        "sortOption": string,
        "LastModified": number,
        "id": string,
        "versionNo": number
    },

    "cardList": {
        "id": versionNo
    },
}

Response:
{
    "requestVersionNo": number,

    "user": {
        "email": string,
        "activated": boolean,
        "isLoggedIn": boolean,
        "rememberMe": boolean,
        "sortOption": string,
        "LastModified": number,
        "id": string,
        "versionNo": number
    },

    "cardToCreate": [
        {
            "belongTo": string,
            "collectedAt": number,
            "context": string,
            "createdAt": number,
            "createdBy": string,
            "detail": string,
            "sourceLang": string,
            "target": string,
            "targetLang": string,
            "translation": string,
            "HasTag": array,
            "collectedBy": string,
            "lastModified": number,
            "id": string,
            "versionNo": number
        },
        {
            ...
        },
    ],
    "cardToUpdate": [
        {
            "belongTo": string,
            "collectedAt": number,
            "context": string,
            "createdAt": number,
            "createdBy": string,
            "detail": string,
            "sourceLang": string,
            "target": string,
            "targetLang": string,
            "translation": string,
            "HasTag": array,
            "collectedBy": string,
            "lastModified": number,
            "id": string,
            "versionNo": number
        },
        {
            ...
        },
    ],
    "cardToDelete": [
        {
            "id": string,
            "versionNo": number
        },
        {
            ...
        },
    ]
}
*/

type CardToDeleteInRes struct {
	ID        bson.ObjectId `json:"id"`
	VersionNo int64         `json:"versionNo"`
}

type ResSignUpOrIn struct {
	User   UserInCommon   `json:"user"`
	Tokens TokensInCommon `json:"tokens"`
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
	Id                string `json:"id"`
	Text              string `json:"text"`
	ChildrenUpdatedAt int64  `json:"childrenUpdatedAt"`
}

type ResDicResults struct {
	// Level: translation/detail/context
	Level   string         `json:"level"`
	Results []DicTextInRes `json:"results"`
}

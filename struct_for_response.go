package testView

import {
	"encoding/json"
	"net/http"
    "reflect"
	"github.com/codegangsta/martini"
    "labix.org/v2/mgo"
}

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
    ID bson.ObjectId
    VersionNo int64
}

type ResSignUpOrIn struct {
    User UserInCommon
    Tokens TokensInCommon
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
    UserId string
    Tokens TokensInCommon
    // No RequestVersionNo here. Tokens are not compared by VersionNo/RequestVersionNo. They are compared directly.
    // No need to have user_id as well. Coz a given set of tokens are only valid with one device and user. If anything goes wrong, simply ask the user to login will do.
}

// Sync user
type ResUser struct {
    User UserInCommon
}

type ResResetPassword struct {
    PasswordIsResetted bool
}

type ResActivation struct {
    IsActivatedThisTime bool
}

type ResDeviceInfo struct {
    DeviceInfo DeviceInfoInCommon
}

// For full card info returned in res
type ResCards struct {
    Cards []CardInCommon
}

type ResSync struct {
    User UserInCommon
    DeviceInfo DeviceInfoInCommon
    CardList []CardInCommon
    CardToDelete []bson.ObjectId
}

type ResFail struct {
    FailedReason string
}

type DicTextInRes {
    Id string
    Text string
    ChildrenUpdatedAt int64
}

type ResDicResults struct {
    // Level: translation/detail/context
    Level string
    Results []DicTextInRes
}
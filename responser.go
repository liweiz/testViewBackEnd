package testView

import (
		"github.com/codegangsta/martini"
		"labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "encoding/json"
        "net/http"
        "string"
        "reflect"
)

/*
Routes:

1 User

1.1 SignUp: 		  	POST 	/users
1.2 SignIn: 		  	GET  	/users/signin
1.3 GetMyAccount:  		GET  	/users/{id}
1.4 ChangeEmail:   		POST 	/users/{id}
1.5 Activate: 	  		POST 	/users/{id}/activation
1.6 ResetPassword: 		POST 	/users/{id}/passwordresetting

Response
1.1.1 	200  ok				user
1.1.2 	409  conflict 		user already exists
1.2.1 	200	 ok				user
1.2.2   401	 unauthorized	reason(incorrect email and/or password)
1.3.1	200  ok				user
1.4.1	200  ok				n/a
1.4.2   409  conflict 		email already in use
1.5.1	show an html page states the success of the activation
1.5.2 	show an error message on the html page
1.6.1	show an html page states the success of the resetting
1.6.2 	show an error message on the html resetting page

*** Access unauthorized related responses	 ***
*** 401  unauthorized	next(please sign in) ***

JSON in request:
1.1 SignUp:
{
    "user": {
    	"email": string
    	"password": string
    }
}

1.5 Activate:
{
    "requestVersionNo": number,

    "user": {
    	"email": string,
    	"activated": boolean,
    	"id": string,
    	"versionNo": number
    }
}

1.6 ResetPassword:
{
    "requestVersionNo": number,

    "user": {
    	"email": string,
    	"activated": boolean,
    	"isLoggedIn": boolean,
    	"id": string,
    	"versionNo": number
    }
}

2 Card

Language options, e.g., Chinese => English, are not shown in URL design.
It is filtered by client side.
2.1 CreateOne: 	  		POST 	/users/{id}/cards
2.2 GetOne: 		  	GET 	/users/{id}/cards/{id}
2.3 DeleteOne: 	  		DELETE 	/users/{id}/cards/{id}
2.4 UpdateOne: 	  		POST 	/users/{id}/cards/{id}
2.5 GetAll: 		  	GET 	/users/{id}/cards

JSON in request:
2.1 CreateOne:
2.4 UpdateOne:
{
    "requestVersionNo": number,
    
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
    }
}

*** Request/Response pair info can be found at: http://stackoverflow.com/questions/16288102/is-it-possible-to-receive-out-of-order-responses-with-http
*** RequestVersionNo in client side's request pool is evaluated by client with the async req/res pair. So no need to have requestVersionNo in response.

Response
2.1.1	200  ok				card
2.1.2 	409  conflict 		conflict detail and card on server that causes the conflict
2.2.1 	200  ok				card
2.2.2 	404  not found		n/a
2.3.1 	200  ok 			n/a
2.3.2 	409  conflict 		conflict detail and card on server that causes the conflict
2.4.1 	200  ok 			card
2.4.2 	409  conflict 		conflict detail and card(s), possiblely to have two cards, on server that causes the conflict
2.5.1 	200  ok 			all cards

3 CardList

This is for sync
3.1 Sync: 		  		POST 	/users/{id}/cardlist

JSON in request:
3.1 Sync:
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
    }
}


Response
3.1.1	200  ok				number of cards to sync(both add/update and delete) and cards
3.1.2 	409  conflict 		conflict detail and card on server that causes the conflict
3.1.3	500  serError 		error detail

4 SearchFromAnyLanguage

LanguagePairs are like Chinese => English.
4.1 StartingWords: 		GET 	/languagepairs/{id}/{id}
4.2 SelectTranslation:	GET 	/languagepairs/{id}/{id}/{id}
4.3 SelectDetail:		GET 	/languagepairs/{id}/{id}/{id}/{id}
4.4 SelectContext:		GET 	/languagepairs/{id}/{id}/{id}/{id}/{id}

Response
4.1.1	200  ok				dictionaryTranslation/dictionaryTarget list
4.2.1 	200  ok 			dictionaryDetail list
4.3.1 	200  ok				dictionaryContext list
4.4.1 	200  ok 			dictionaryContext


Response in summary			

200  ok
		user 		1.1.1, 1.2.1, 1.3.1
		n/a 		1.4.1, 2.3.1
		card 		2.1.1, 2.2.1, 2.4.1
		all cards 	2.5.1
		number of cards to sync(both add/update and delete) and cards  	3.1.1
		dictionaryTranslation/dictionaryTarget list 					4.1.1
		dictionaryDetail list 											4.2.1
		dictionaryContext list 											4.3.1
		dictionaryContext 												4.4.1

409  conflict
		user already exists												1.1.2
		email already in use 											1.4.2
		conflict detail and card on server that causes the conflict 	2.1.2, 2.3.2, 3.1.2
		conflict detail and card(s), possiblely to have two cards, on server that causes the conflict  2.4.2

404  not found
		n/a 		2.2.2

401  unauthorized
		next(please sign in)
		reason(incorrect email and/or password)   1.2.2

1.5.1	show an html page states the success of the activation
1.5.2 	show an error message on the html page
1.6.1	show an html page states the success of the resetting
1.6.2 	show an error message on the html resetting page
*/

// Decisions on requests
const (
	// 200 ok
	OkOnly											= iota
	OkAndUser
	OkAndDeviceInfo
	OkAndCards
	OkAndAllCards
	// number of cards to sync(both add/update and delete) and cards
	OkAndSync
	// dictionaryTranslation/dictionaryTarget list
	OkAndDictWordsList
	OkAndDictDetailList
	OkAndDictContextList

	// 409 conflict
	ConflictUserAlreadyExists
	ConflictEmailAlreadyInUse
	// For new card only
	ConflictCardAlreadyExists
	// conflict detail and card(s), possiblely to have two cards, on server that causes the conflict
	ConflictDetailAndCardsOverwriteClient
	ConflictDetailAndCardsOverwriteDB
	ConflictDetailAndCardsCreateAnotherInDB\
	// For DeviceInfo and User
	OverWriteClient
	OverWriteServer
	OverWriteServerEmail

	// 404 not found
	NotFoundOnly

	// 405 method not allowed
	MethodNotAllowed

	// 500 Internal Server Error
	InternalServerError

	// 401 unauthorized
	UnauthorizedWithSignInReq
	// incorrect email and/or password
	UnauthorizedWithDetail
	
	// serve html
	ActivationDone
	ActivationNotDone
	PasswordResettingDone
	PasswordResettingNotDone

	// potential error 
	HasErr

)

// According to the decision, handle conflict, get the data needed from db and write the response.
func GenerateResponse(decision int, message string, rw martini.ResponseWriter, structPartInRes interface{}) {
	rw := martini.NewResponseWriter()
	rw.Header().Set("Content-Type", "application/json")
	switch decision {
		case OkOnly:
			rw.WriteHeader(StatusOK)
		case OkAndUser:
			var user User
			
			EncodeAndResponse(err, user)
		case OkAndCard:
			var card Card
			err := db.C("cards").Find(criteria).One(&card)
			EncodeAndResponse(err, card)
		case OkAndAllCards:
			var cards []Card{}
			err := db.C("cards").Find(criteria).All(&cards)
			EncodeAndResponse(err, cards)
		case OkAndSync:
			// There are one slice for add/update and another for delete
			// Iterate the add/update one first

			// And iterate the delete one

		case OkAndDictWordsList:
		case OkAndDictDetailList:
		case OkAndDictContextList:
		case ConflictUserAlreadyExists:
		case ConflictEmailAlreadyInUse:
		case ConflictDetailAndCardsOverwriteClient:
		case ConflictDetailAndCardsOverwriteDB:
		case ConflictDetailAndCardsCreateAnotherInDB:
		case NotFoundOnly:
		case MethodNotAllowed:
		case InternalServerError:
		case UnauthorizedWithSignInReq:
		case UnauthorizedWithDetail:
	}
}





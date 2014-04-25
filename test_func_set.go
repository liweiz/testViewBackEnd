package main

import (
	// "me/testView/handlers"
	"fmt"
	"net/http"
)

type funcForTestStep func() (*http.Request, int)

func GetFuncForTestStep(m *MyMartini, p *publicDataSet, u1 *userX, u2 *userX) map[string]funcForTestStep {
	return map[string]funcForTestStep{
		"SignUp": func() (*http.Request, int) {
			fmt.Println("================================= SignUp =================================")
			reqBodyStruct := p.GetSignUpReqBodyStruct()
			return p.TestSignUp(m, reqBodyStruct), 0
		},
		"SignIn": func() (*http.Request, int) {
			fmt.Println("================================= SignIn =================================")
			reqBodyStruct := p.GetSignInReqBodyStruct()
			return p.TestSignIn(m, reqBodyStruct), 1
		},
		"RenewTokens": func() (*http.Request, int) {
			fmt.Println("================================= RenewTokens =================================")
			reqBodyStruct := p.GetRenewTokensReqBodyStruct()
			return p.TestRenewTokens(m, reqBodyStruct), 2
		},
		"ActivationEmail": func() (*http.Request, int) {
			fmt.Println("================================= ActivationEmail =================================")
			return p.TestActivationEmail(m), 3
		},
		"ClickActivationLink": func() (*http.Request, int) {
			fmt.Println("================================= ClickActivationLink =================================")
			return p.TestClickActivationLink(m), 4
		},
		"PasswordResettingEmailByToken": func() (*http.Request, int) {
			fmt.Println("================================= PasswordResettingEmailByToken =================================")
			return p.TestPasswordResettingEmailByToken(m), 5
		},
		"PasswordResettingEmailByEmail": func() (*http.Request, int) {
			fmt.Println("================================= PasswordResettingEmailByEmail =================================")
			reqBodyStruct := p.GetForgotPwdSendEmailToResetReqBodyStruct()
			return p.TestPasswordResettingEmailByEmail(m, reqBodyStruct), 6
		},
		"ChangePassword": func() (*http.Request, int) {
			fmt.Println("================================= ChangePassword =================================")
			reqBodyStruct := p.GetResetPasswordReqBodyStruct()
			return p.TestChangePassword(m, reqBodyStruct), 7
		},
		"NewDeviceInfo": func() (*http.Request, int) {
			fmt.Println("================================= NewDeviceInfo =================================")
			reqBodyStruct := p.GetDeviceInfoReqBodyStruct()
			return p.TestNewDeviceInfo(m, reqBodyStruct), 8
		},
		"UpdateDeviceInfo": func() (*http.Request, int) {
			fmt.Println("================================= UpdateDeviceInfo =================================")
			reqBodyStruct := p.GetDeviceInfoReqBodyStruct()
			return p.TestUpdateDeviceInfo(m, reqBodyStruct), 9
		},
		"NewCardV1": func() (*http.Request, int) {
			fmt.Println("================================= NewCard V1 =================================")
			p.SetPublicDetail(1)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestNewCard(m, reqBodyStruct), 10
		}, // new card
		"NewCardV2": func() (*http.Request, int) {
			fmt.Println("================================= NewCard V2 =================================")
			p.SetPublicDetail(2)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestNewCard(m, reqBodyStruct), 11
		}, // new card
		"NewCardV3": func() (*http.Request, int) {
			fmt.Println("================================= NewCard V3 =================================")
			p.SetPublicDetail(3)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestNewCard(m, reqBodyStruct), 12
		}, // new card
		"UpdateCardV1": func() (*http.Request, int) {
			fmt.Println("================================= UpdateCard V1  =================================")
			p.SetPublicDetail(1)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestUpdateCard(m, reqBodyStruct), 13
		}, // update card
		"UpdateCardV2": func() (*http.Request, int) {
			fmt.Println("================================= UpdateCard V2  =================================")
			p.SetPublicDetail(2)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestUpdateCard(m, reqBodyStruct), 14
		}, // update card
		"UpdateCardV3": func() (*http.Request, int) {
			fmt.Println("================================= UpdateCard V3  =================================")
			p.SetPublicDetail(3)
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestUpdateCard(m, reqBodyStruct), 15
		}, // update card
		"DeleteCard": func() (*http.Request, int) {
			fmt.Println("================================= DeleteCard =================================")
			reqBodyStruct := p.GetCardReqBodyStruct()
			return p.TestDeleteCard(m, reqBodyStruct), 16
		}, // delete card
		"SyncEmptyCard": func() (*http.Request, int) {
			fmt.Println("================================= SyncEmptyCard =================================")
			reqBodyStruct := p.GetSyncReqBodyStructEmptyCard()
			return p.TestSync(m, reqBodyStruct), 17
		}, // sync
		"SyncEmptyCardNewDevice": func() (*http.Request, int) {
			fmt.Println("================================= SyncEmptyCardNewDevice =================================")
			reqBodyStruct := p.GetSyncReqBodyStructEmptyCardNewDevice()
			return p.TestSync(m, reqBodyStruct), 18
		}, // sync
		"SyncWithCards": func() (*http.Request, int) {
			fmt.Println("================================= SyncWithCards =================================")
			reqBodyStruct := p.GetSyncReqBodyStructWithCards()
			return p.TestSync(m, reqBodyStruct), 19
		}, // sync
		"SyncEmptyCardNewDevice2": func() (*http.Request, int) {
			fmt.Println("================================= SyncEmptyCardNewDevice2 =================================")
			reqBodyStruct := p.GetSyncReqBodyStructEmptyCardNewDevice2()
			return p.TestSync(m, reqBodyStruct), 20
		}, // sync
		"assign tokens to user1": func() (*http.Request, int) {
			fmt.Println("================================= assign tokens to user1 =================================")
			u1.SetUserTokens(p)
			return nil, 100
		},
		"assign tokens to user2": func() (*http.Request, int) {
			fmt.Println("================================= assign tokens to user2 =================================")
			u2.SetUserTokens(p)
			return nil, 101
		},
		"assign user1 tokens to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user1 tokens to public =================================")
			// p.Uuid takes them to the right device
			p.SetPublicTokens(u1)
			return nil, 102
		},
		"assign user2 tokens to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user2 tokens to public =================================")
			p.SetPublicTokens(u2)
			return nil, 103
		},
		"assign device1 uuid to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 uuid to public =================================")
			p.SetPublicUuid(u1.Devices[0].DeviceUuid)
			return nil, 104
		},
		"assign device2 uuid to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device2 uuid to public =================================")
			p.SetPublicUuid(u1.Devices[1].DeviceUuid)
			return nil, 105
		},
		"assign user2 email to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user2 email to public =================================")
			p.SetPublicEmail(u2.Email)
			return nil, 106
		},
		"assign user1 email to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user1 email to public =================================")
			p.SetPublicEmail(u1.Email)
			return nil, 107
		},
		"assign user1 activation url code to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user1 activation url code to public =================================")
			p.GetActivationUrlCodeFromDb(u1)
			return nil, 108
		},
		"assign user2 activation url code to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user2 activation url code to public =================================")
			p.GetActivationUrlCodeFromDb(u2)
			return nil, 109
		},
		"assign user1 pwd resetting url code to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user1 pwd resetting url code to public =================================")
			p.GetPasswordResettingUrlCodeFromDb(u1)
			return nil, 110
		},
		"assign user2 pwd resetting url code to public": func() (*http.Request, int) {
			fmt.Println("================================= assign user2 pwd resetting url code to public =================================")
			p.GetPasswordResettingUrlCodeFromDb(u2)
			return nil, 111
		},
		"assign device1 2nd reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 2nd reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[1]
			return nil, 112
		},
		"assign device1 3rd reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 3rd reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[2]
			return nil, 113
		},
		"assign device1 4th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 4th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[3]
			return nil, 114
		},
		"assign device1 5th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 5th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[4]
			return nil, 115
		},
		"assign sort option '1' to public": func() (*http.Request, int) {
			fmt.Println("================================= assign sort option '1' to public =================================")
			p.SortOption = "1"
			return nil, 116
		},
		"assign sort option '2' to public": func() (*http.Request, int) {
			fmt.Println("================================= assign sort option '2' to public =================================")
			p.SortOption = "2"
			return nil, 117
		},
		"assign card id original ver no equal to db's to public": func() (*http.Request, int) {
			fmt.Println("================================= assign card id original ver no equal to db's to public =================================")
			p.GetCardVerNoFromDb()
			return nil, 118
		},
		// These two should be used after "assign card id original ver no equal to db's to public", or there is no way to get the p.CardIdOriginalVerNo in the first place.
		"assign card id original ver no less than db's to public": func() (*http.Request, int) {
			fmt.Println("================================= assign card id original ver no less than db's to public =================================")
			p.CardIdOriginalVerNo = p.CardIdOriginalVerNo - 1
			return nil, 119
		},
		"assign card id original ver no greater than db's to public": func() (*http.Request, int) {
			fmt.Println("================================= assign card id original ver no greater than db's to public =================================")
			p.CardIdOriginalVerNo = p.CardIdOriginalVerNo + 1
			return nil, 120
		},
		"assign device1 6th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 6th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[5]
			return nil, 121
		},
		"assign device1 7th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 7th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[6]
			return nil, 122
		},
		"assign device1 8th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 8th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[7]
			return nil, 123
		},
		"assign device1 9th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 9th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[8]
			return nil, 124
		},
		"assign device1 10th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 10th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[9]
			return nil, 125
		},
		"assign device1 11th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 11th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[10]
			return nil, 126
		},
		"assign device1 12th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 12th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[11]
			return nil, 127
		},
		"assign device1 13th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 13th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[12]
			return nil, 128
		},
		"assign device1 14th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 14th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[13]
			return nil, 129
		},
		"assign device1 15th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 15th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[14]
			return nil, 130
		},
		"assign device1 16th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 16th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[15]
			return nil, 131
		},
		"assign device1 17th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 17th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[16]
			return nil, 132
		},
		"assign device1 18th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 18th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[17]
			return nil, 133
		},
		"assign device1 19th reqId to public": func() (*http.Request, int) {
			fmt.Println("================================= assign device1 19th reqId to public =================================")
			p.ReqId = u1.Devices[0].ReqId[18]
			return nil, 134
		},
		"remove all cards in db": func() (*http.Request, int) {
			fmt.Println("================================= remove all cards in db =================================")
			p.ClearAllCardsInDb()
			return nil, 135
		},
		"insert cards in db for sync test": func() (*http.Request, int) {
			fmt.Println("================================= insert cards in db for sync test =================================")
			p.SetSyncTestCardsInDb()
			return nil, 136
		},
	}
}

package main

import (
	// "me/testView/handlers"
	// "fmt"
	// "net/http"
	// "net/http/httptest"
	"testing"
)

func TestSteps(t *testing.T) {
	m := My()
	// Start point
	p := &publicDataSet{}

	d := GetTestDataSet()
	// Setup user1
	u1 := &userX{}
	u1d1 := deviceX{}
	u1d2 := deviceX{}
	u1.Devices = append(u1.Devices, u1d1, u1d2)
	u1.UserNo = 1
	u1.Email = d["email"][0]
	u1.Password = d["password"][0]
	u1.Devices[0].DeviceUuid = d["uuid"][0]
	u1.Devices[0].DeviceNo = 1
	u1.Devices[1].DeviceUuid = d["uuid"][1]
	u1.Devices[1].DeviceNo = 2
	// Use the same set of reqId to test if same reqId from different device could go wrong.
	u1.Devices[0].ReqId = d["reqId1"]
	u1.Devices[1].ReqId = d["reqId1"]
	// Setup user2
	u2 := &userX{}
	u2d1 := deviceX{}
	u2.Devices = append(u2.Devices, u2d1)
	u2.UserNo = 2
	u2.Email = d["email"][1]
	u2.Password = d["password"][1]
	u2.Devices[0].DeviceUuid = d["uuid"][0]
	u2.Devices[0].DeviceNo = 1
	u2.Devices[0].ReqId = d["reqId1"]

	f := GetFuncForTestStep(m, p, u1, u2)
	// setup start point
	p.Email = u1.Email
	p.Password = u1.Password
	p.Uuid = u1.Devices[0].DeviceUuid
	p.ReqId = u1.Devices[0].ReqId[0]
	p.SortOption = "1"
	fTestFlow := []funcForTestStep{
		// user1 signs up, signs in, renew tokens
		f["SignUp"],                        // user1 signs up on device1
		f["assign tokens to user1"],        // assign tokens to user1device1
		f["SignUp"],                        // try sign up with the dulpicated email
		f["SignIn"],                        // user1 signs in on device1
		f["assign tokens to user1"],        // assign tokens to user1device1
		f["assign device2 uuid to public"], // change uuid to device2's
		f["SignIn"],                        // user1 signs in on device2
		f["assign tokens to user1"],        // assign tokens to user1device2
		f["assign device1 uuid to public"], // change uuid to device1's
		f["assign user1 tokens to public"], // change tokens to user1device1's
		f["RenewTokens"],                   // renew tokens for user1device1
		f["assign tokens to user1"],        // assign tokens to user1device1

		// activate user1
		f["ActivationEmail"],                            // send activation email for user1
		f["assign user1 activation url code to public"], // get and assign user1 activation url code to public
		f["ClickActivationLink"],                        // activate user1

		// change user1's password
		// f["PasswordResettingEmailByToken"],                 // get user1 pwd resetting email by token
		// f["assign user1 pwd resetting url code to public"], // get and assign user1 pwd resetting url code to public
		// f["ChangePassword"],                                // change user1 password through url
		// f["PasswordResettingEmailByEmail"],                 // get user1 pwd resetting email by email
		// f["assign user1 pwd resetting url code to public"], // get and assign user1 pwd resetting url code to public
		// f["ChangePassword"],                                // change user1 password through url

		// create and update new deviceInfo, pls see info architecture.xlsx sheet10
		// starting point
		f["NewDeviceInfo"], // create a device info for device1
		f["assign device1 2nd reqId to public"],
		f["assign sort option '2' to public"],
		f["UpdateDeviceInfo"], // update a device info for device1

		/////////////////////// single card CRUD BEGINS /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// first card, this no is from Excel file column E
		// 2
		// f["assign device1 3rd reqId to public"],
		// f["NewCardV1"], // this is the original card since it's the first one created
		// // 1
		// f["assign card id original ver no equal to db's to public"],
		// f["assign device1 4th reqId to public"],
		// f["NewCardV1"], // uniqueness test, but it has to be with a different request id:) the request id is not supposed to be added to db since it can not pass the uniqueness check, which generates an err.
		// // 7
		// f["assign device1 5th reqId to public"],
		// f["UpdateCardV1"], // update original card with same content, should not be able to pass uniqueness check
		// // 8
		// f["assign device1 6th reqId to public"],
		// f["UpdateCardV2"], // update original card with different content and same version no
		// // 10
		// f["assign device1 7th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no less than db's to public"],
		// f["UpdateCardV1"], // update original card with different content and version no less than db's, should be ConflictCreateAnotherInDB
		// // 9
		// f["assign device1 8th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no greater than db's to public"],
		// f["UpdateCardV3"], // should be successful with NoConflictOverwriteDB
		// // 12
		// f["assign device1 9th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["DeleteCard"], // should be successful with overwriting client to make sure no inconsistence
		// // 11
		// f["assign device1 10th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["DeleteCard"], // should be successful with NoConflictOverwriteDB
		// // 3
		// f["assign device1 11th reqId to public"],
		// f["UpdateCardV1"], // should not pass uniqueness check
		// // 4
		// f["assign device1 12th reqId to public"],
		// f["UpdateCardV2"], // should be ConflictCreateAnotherInDB
		// // 6
		// f["assign device1 13th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no less than db's to public"],
		// f["UpdateCardV3"], // should be ConflictCreateAnotherInDB
		// // remove all cards to have fresh restart, update once to make the version no equal ot 2
		// f["remove all cards in db"],
		// f["assign device1 14th reqId to public"],
		// f["NewCardV1"],
		// f["assign device1 15th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["UpdateCardV2"],
		// // 14
		// f["assign device1 16th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no less than db's to public"],
		// f["DeleteCard"], // keep the card and overwrite client's, this provides another chance for user to make decision based on updated content
		// // 13
		// f["assign device1 17th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no greater than db's to public"],
		// f["DeleteCard"], // overwrite client to make sure no inconsistence
		// // 5
		// f["assign device1 18th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["DeleteCard"], // actually delete the card to create the situation for 5
		// f["assign device1 19th reqId to public"],
		// f["assign card id original ver no equal to db's to public"],
		// f["assign card id original ver no greater than db's to public"],
		// f["UpdateCardV1"],
		/////////////////////// single card CRUD ENDS /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	}
	RunOperationFlow(fTestFlow, m, p)

}

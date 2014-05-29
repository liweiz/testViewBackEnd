package main

import (
	// "bytes"
	"encoding/json"
	//"github.com/go-martini/martini"
	"fmt"
	// "labix.org/v2/mgo"
	// "labix.org/v2/mgo/bson"
	"me/testViewPKG/testView"
	"net/http"
	"net/http/httptest"
	"os"
)

// s stands for any sample data slice.
// s[0]: user1 or device1 data
// s[1]: user2 or device2 data
// s[3]: random data that is not correct

func GetTestDataSet() map[string][]string {
	return map[string][]string{
		"email":    []string{"matt.z.lw@gmail.com", "remolet.z@gmail.com", "aUser"},
		"password": []string{"aA1~aA1~", "aA1~aA1~!", "xX~xX~"},
		"uuid":     []string{"aaa", "bbb", "ccc"},
		"reqId1":   []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "a10", "a11", "a12", "a13", "a14", "a15", "a16", "a17", "a18", "a19", "a20", "a21", "a22", "a23", "a24", "a25", "a26", "a27", "a28", "a29", "a30", "a31", "a32", "a33", "a34", "a35", "a36", "a37", "a38", "a39"},
		"reqId2":   []string{"b1", "b2", "b3", "b4", "b5", "b6"},
	}
}

func RunOperationFlow(OperationFlow []funcForTestStep, m *MyMartini, p *publicDataSet) {
	for _, f := range OperationFlow {
		req, i := f()
		if i < 100 {
			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)
			fmt.Println("code:", w.Code)
			os.Stdout.Write(w.Body.Bytes())
			if w.Code == http.StatusOK {
				d := w.Body.Bytes()
				var err error
				// DeviceInfo
				if i == 9 {
					i = 8
				}
				// Create card
				if i == 11 || i == 12 {
					i = 10
				}
				// Update card
				if i == 14 || i == 15 {
					i = 13
				}
				switch i {
				case 0:
					r := &testView.ResSignUpOrIn{}
					err = json.Unmarshal(d, r)
					p.UserId = r.User.Id.Hex()
					p.AccessToken = r.Tokens.AccessToken
					p.RefreshToken = r.Tokens.RefreshToken
				case 1:
					r := &testView.ResSignUpOrIn{}
					err = json.Unmarshal(d, r)
					p.UserId = r.User.Id.Hex()
					p.AccessToken = r.Tokens.AccessToken
					p.RefreshToken = r.Tokens.RefreshToken
				case 2:
					r := &testView.ResTokensOnly{}
					err = json.Unmarshal(d, r)
					p.AccessToken = r.Tokens.AccessToken
					p.RefreshToken = r.Tokens.RefreshToken
				case 8:
					r := &testView.ResDeviceInfo{}
					err = json.Unmarshal(d, r)
					p.DeviceInfoId = r.DeviceInfo.Id.Hex()
				case 9:
					r := &testView.ResDeviceInfo{}
					err = json.Unmarshal(d, r)
					p.DeviceInfoId = r.DeviceInfo.Id.Hex()
				case 10:
					r := &testView.ResCards{}
					err = json.Unmarshal(d, r)
					if len(r.Cards) > 0 {
						p.CardIdOriginal = r.Cards[0].Id.Hex()
					}
				case 13:
					r := &testView.ResCards{}
					err = json.Unmarshal(d, r)
					if len(r.Cards) == 2 {
						if r.Cards[0].Id.Hex() == p.CardIdOriginal {
							p.CardIdDerived = r.Cards[1].Id.Hex()
						} else {
							p.CardIdDerived = r.Cards[0].Id.Hex()
						}
					}
				case 21:
					r := &testView.ResDicResultsText{}
					err = json.Unmarshal(d, r)
					if len(r.Results) > 0 {
						p.LastIdInSearch = r.Results[len(r.Results)-1].Id.Hex()
					}
					p.ParentIdInSearch = r.TopLevelTextId.Hex()
				case 22:
					r := &testView.ResDicResultsId{}
					err = json.Unmarshal(d, r)
					if len(r.Results) > 0 {
						p.LastIdInSearch = r.Results[len(r.Results)-1].Id.Hex()
					}
				default:
					//
				}
				if err != nil {
					fmt.Println("err: ", err.Error())
				}
			}
		} else {

		}

	}
}

// os.Stdout.Write(reqBody)

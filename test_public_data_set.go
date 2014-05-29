package main

import (
	// "me/testView/handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"me/testViewPKG/testView"
	"net/http"
	// "net/http/httputil"
	"os"
	"time"
)

type publicDataSet struct {
	Email                    string
	Password                 string
	UserId                   string
	AccessToken              string
	RefreshToken             string
	ActivationUrlCode        string
	PasswordResettingUrlCode string
	Uuid                     string
	DeviceInfoId             string
	CardIdOriginal           string
	CardIdOriginalVerNo      int64
	CardIdDerived            string
	ReqId                    string
	SortOption               string
	Detail                   string
	SyncTestCardId1          string
	SyncTestCardId2          string
	SyncTestCardId3          string
	SyncTestCardId4          string
	SyncTestCardId5          string
	SyncTestCardId6          string
	SyncTestCardId7          string
	SyncTestCardId8          string
	TextToSearch             string
	ParentIdInSearch         string
	LastIdInSearch           string
	SearchSortOption         string
	SearchAscending          bool
}

///////////TEST FUNC/////////////
func (p *publicDataSet) TestSignUp(m *MyMartini, reqBodyStruct *testView.ReqSignUpOrIn) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users"
	m.Post(url, testView.RequestPreprocessor(testView.SignUp), testView.ProcessedResponseGenerator(testView.SignUp, false))
	req, _ = http.NewRequest("POST", url, body)
	return
}

func (p *publicDataSet) TestSignIn(m *MyMartini, reqBodyStruct *testView.ReqSignUpOrIn) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/signin"
	m.Post(url, testView.GateKeeper(), testView.RequestPreprocessor(testView.SignIn), testView.ProcessedResponseGenerator(testView.SignIn, false))
	req, _ = http.NewRequest("POST", url, body)
	req.SetBasicAuth(p.Email, p.Password)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestRenewTokens(m *MyMartini, reqBodyStruct *testView.ReqRenewTokens) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/tokens"
	m.Post("/users/:user_id/tokens", testView.GateKeeperExchange(), testView.RequestPreprocessor(testView.RenewTokens), testView.ProcessedResponseGenerator(testView.RenewTokens, false))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestActivationEmail(m *MyMartini) (req *http.Request) {
	url := "/users/" + p.UserId + "/activation"
	m.Get("/users/:user_id/activation", testView.GateKeeper(), testView.EmailSender(testView.EmailForActivation))
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestClickActivationLink(m *MyMartini) (req *http.Request) {
	url := "/users/" + p.UserId + "/activation/" + p.ActivationUrlCode
	m.Get("/users/:user_id/activation/:activation_code", testView.WebPageServer(testView.PageForActivation))
	req, _ = http.NewRequest("GET", url, nil)
	return
}

// This is for reset password when logged in. So gateKeeper needed.
func (p *publicDataSet) TestPasswordResettingEmailByToken(m *MyMartini) (req *http.Request) {
	url := "/users/" + p.UserId + "/password"
	m.Get("/users/:user_id/password", testView.GateKeeper(), testView.EmailSender(testView.EmailForPasswordResetting))
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

// This is for forgot-password. No gateKeeper needed, but user has to provide the email to verify if the account exsits. If yes, send the email and let user know. Otherwise, tell the user the account does not exist.
func (p *publicDataSet) TestPasswordResettingEmailByEmail(m *MyMartini, reqBodyStruct *testView.ReqForgotPassword) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/forgotpassword"
	m.Post("/users/forgotpassword", testView.RequestPreprocessor(testView.ForgotPassword), testView.ProcessedResponseGenerator(testView.ForgotPassword, false))
	m.Get("/assets/css/bootstrap.min.css", testView.AssetsServer(testView.BootstrapCssMin))
	m.Get("/assets/js/bootstrap.min.js", testView.AssetsServer(testView.BootstrapJsMin))
	m.Get("/assets/css/account_activation.css", testView.AssetsServer(testView.CssPageForActivation))
	m.Get("/assets/js/account_activation.js", testView.AssetsServer(testView.JsPageForActivation))
	m.Get("/assets/css/password_resetting.css", testView.AssetsServer(testView.CssPageForPasswordResetting))
	m.Get("/assets/js/password_resetting.js", testView.AssetsServer(testView.JsPageForPasswordResetting))
	req, _ = http.NewRequest("POST", url, body)
	return
}

// Reset password from the uniqueUrl
func (p *publicDataSet) TestChangePassword(m *MyMartini, reqBodyStruct *testView.ReqResetPassword) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/password/" + p.PasswordResettingUrlCode
	m.Post("/users/:user_id/password/:password_resetting_code", testView.UrlCodeChecker(), testView.RequestPreprocessor(testView.PasswordResetting), testView.ProcessedResponseGenerator(testView.PasswordResetting, false))
	req, _ = http.NewRequest("POST", url, body)
	return
}

func (p *publicDataSet) TestNewDeviceInfo(m *MyMartini, reqBodyStruct *testView.ReqDeviceInfo) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/deviceinfos"
	// Use this to simulate use token with wrong p.UserId.
	// url := "/users/53498c7abd7e550ce4000002/deviceinfo"
	m.Post("/users/:user_id/deviceinfos", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.NewDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewDeviceInfo, true))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestUpdateDeviceInfo(m *MyMartini, reqBodyStruct *testView.ReqDeviceInfo) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/deviceinfos/" + p.DeviceInfoId
	m.Post("/users/:user_id/deviceinfos/:device_id", testView.GateKeeper(), testView.RequestPreprocessor(testView.OneDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneDeviceInfo, true))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestNewCard(m *MyMartini, reqBodyStruct *testView.ReqCard) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/cards"
	m.Post("/users/:user_id/cards", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.NewCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewCard, true))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestUpdateCard(m *MyMartini, reqBodyStruct *testView.ReqCard) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/cards/" + p.CardIdOriginal
	m.Post("/users/:user_id/cards/:card_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestDeleteCard(m *MyMartini, reqBodyStruct *testView.ReqCard) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/cards/" + p.CardIdOriginal
	m.Delete("/users/:user_id/cards/:card_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	req, _ = http.NewRequest("DELETE", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	return
}

func (p *publicDataSet) TestSync(m *MyMartini, reqBodyStruct *testView.ReqSync) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/users/" + p.UserId + "/sync"
	m.Post("/users/:user_id/sync", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.Sync), testView.ProcessedResponseGenerator(testView.Sync, false))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	os.Stdout.Write(reqBody)
	return
}

func (p *publicDataSet) TestDicSearchByText(m *MyMartini, reqBodyStruct *testView.ReqDicText) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/dic/7/37/text/" + p.UserId
	m.Post("/dic/:source_lang_code/:target_lang_code/text/:user_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.DicTextSearcher())
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	os.Stdout.Write(reqBody)
	return
}

func (p *publicDataSet) TestDicSearchById(m *MyMartini, reqBodyStruct *testView.ReqDicId) (req *http.Request) {
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	body := bytes.NewReader(reqBody)
	url := "/dic/7/37/id/" + p.UserId
	m.Post("/dic/:source_lang_code/:target_lang_code/id/:user_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.DicTextSearcher())
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)
	req.Header.Add("X-REMOLET-DEVICE-ID", p.Uuid)
	os.Stdout.Write(reqBody)
	return
}

///////////DATA MODIFICATION/////////////
func (p *publicDataSet) SetPublicTokens(u *userX) {
	fmt.Println("From user: ", u.UserNo)
	if u.Devices[0].DeviceUuid == p.Uuid {
		p.AccessToken = u.Devices[0].AccessToken
		p.RefreshToken = u.Devices[0].RefreshToken
		fmt.Println("From device: ", u.Devices[0].DeviceNo)
	} else {
		p.AccessToken = u.Devices[1].AccessToken
		p.RefreshToken = u.Devices[1].RefreshToken
		fmt.Println("From device: ", u.Devices[1].DeviceNo)
	}
}

func (p *publicDataSet) SetPublicUuid(uuid string) {
	fmt.Println("To Public uuid")
	p.Uuid = uuid
}

func (p *publicDataSet) SetPublicEmail(email string) {
	fmt.Println("To Public email")
	p.Email = email
}

func (p *publicDataSet) SetPublicDetail(detailNo int) {
	fmt.Println("To Public detail: ", detailNo)
	switch detailNo {
	case 1:
		p.Detail = "直译为“移情作用”，在中文中不易理解。"
	case 2:
		p.Detail = "直译为“移情作用”。wiki中有详解，却不易理解。"
	case 3:
		p.Detail = "直译为“移情作用”。"
	case 4:
		p.Detail = "a直译为“移情作用”，在中文中不易理解。"
	case 5:
		p.Detail = "a直译为“移情作用”。wiki中有详解，却不易理解。"
	case 6:
		p.Detail = "a直译为“移情作用”。"
	}
}

///////////Request Body Struct/////////////
func (p *publicDataSet) GetSignUpReqBodyStruct() *testView.ReqSignUpOrIn {
	return &testView.ReqSignUpOrIn{
		p.Email,
		p.Password,
		p.Uuid}
}

func (p *publicDataSet) GetSignInReqBodyStruct() *testView.ReqSignUpOrIn {
	return &testView.ReqSignUpOrIn{
		"",
		"",
		p.Uuid}
}

func (p *publicDataSet) GetRenewTokensReqBodyStruct() *testView.ReqRenewTokens {
	return &testView.ReqRenewTokens{
		DeviceUUID: p.Uuid,
		Tokens:     testView.TokensInCommon{p.AccessToken, p.RefreshToken}}
}

func (p *publicDataSet) GetForgotPwdSendEmailToResetReqBodyStruct() *testView.ReqForgotPassword {
	return &testView.ReqForgotPassword{p.Email}
}

func (p *publicDataSet) GetResetPasswordReqBodyStruct() *testView.ReqResetPassword {
	return &testView.ReqResetPassword{p.Password}
}

// New and update have the same body struct.
func (p *publicDataSet) GetDeviceInfoReqBodyStruct() *testView.ReqDeviceInfo {
	return &testView.ReqDeviceInfo{
		RequestId:  p.ReqId,
		DeviceUUID: p.Uuid,
		DeviceInfo: testView.DeviceInfoInCommonNew{
			DeviceUUID: p.Uuid,
			SourceLang: "English",
			TargetLang: "简体中文",
			SortOption: p.SortOption,
			IsLoggedIn: true,
			RememberMe: true}}
}

func (p *publicDataSet) GetCardReqBodyStruct() *testView.ReqCard {
	fmt.Println("CardIdOriginalVerNo: ", p.CardIdOriginalVerNo)
	return &testView.ReqCard{
		RequestId:  p.ReqId,
		DeviceUUID: p.Uuid,
		Card: testView.CardInCommonNew{
			Context:     "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
			Target:      "empathy",
			Translation: "感同身受",
			Detail:      p.Detail,
			SourceLang:  "English",
			TargetLang:  "简体中文"},
		CardVersionNo: p.CardIdOriginalVerNo}
}

func (p *publicDataSet) GetSyncReqBodyStructEmptyCard() *testView.ReqSync {
	return &testView.ReqSync{
		DeviceUUID: p.Uuid}
}

func (p *publicDataSet) GetSyncReqBodyStructEmptyCardNewDevice() *testView.ReqSync {
	return &testView.ReqSync{
		DeviceUUID: "p.Uuid"}
}

func (p *publicDataSet) GetSyncReqBodyStructEmptyCardNewDevice2() *testView.ReqSync {
	return &testView.ReqSync{
		DeviceUUID: "p.Uuid2"}
}

func (p *publicDataSet) SetSyncTestCardsInDb() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	d1 := bson.NewObjectId()
	p.SyncTestCardId1 = d1.Hex()
	SyncTestCard1 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "11111",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d1,
		"versionNo":            2,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            false}
	d2 := bson.NewObjectId()
	p.SyncTestCardId2 = d2.Hex()
	SyncTestCard2 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "22222",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d2,
		"versionNo":            1,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            false}
	d3 := bson.NewObjectId()
	p.SyncTestCardId3 = d3.Hex()
	SyncTestCard3 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "33333",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d3,
		"versionNo":            1,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            false}
	d4 := bson.NewObjectId()
	p.SyncTestCardId4 = d4.Hex()
	SyncTestCard4 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "44444",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d4,
		"versionNo":            1,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            false}
	d5 := bson.NewObjectId()
	p.SyncTestCardId5 = d5.Hex()
	SyncTestCard5 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "55555",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d5,
		"versionNo":            2,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            true}
	d6 := bson.NewObjectId()
	p.SyncTestCardId6 = d6.Hex()
	SyncTestCard6 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "66666",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d6,
		"versionNo":            1,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            true}
	d7 := bson.NewObjectId()
	p.SyncTestCardId7 = d7.Hex()
	SyncTestCard7 := bson.M{
		"context":              "As designers, we must not forget that we design for the people. We must gain empathy and ride on the arc of modern design.",
		"target":               "empathy",
		"translation":          "感同身受",
		"detail":               "77777",
		"sourceLang":           "English",
		"targetLang":           "简体中文",
		"_id":                  d7,
		"versionNo":            1,
		"belongTo":             bson.ObjectIdHex(p.UserId),
		"lastModified":         time.Now().UnixNano(),
		"collectedAt":          time.Now().UnixNano(),
		"contextDicTextId":     "",
		"detailDicTextId":      "",
		"targetDicTextId":      "",
		"translationDicTextId": "",
		"isDeleted":            true}
	_ = s.DB("mylang").C("cards").Insert(SyncTestCard1, SyncTestCard2, SyncTestCard3, SyncTestCard4, SyncTestCard5, SyncTestCard6, SyncTestCard7)
	defer s.Close()
	fmt.Println("p.SyncTestCardId7: ", p.SyncTestCardId7)
}

func (p *publicDataSet) GetSyncReqBodyStructWithCards() *testView.ReqSync {
	return &testView.ReqSync{
		DeviceUUID: p.Uuid,
		CardList: []testView.CardsVerListElement{
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId1), 1},
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId2), 1},
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId3), 2},
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId5), 1},
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId6), 1},
			testView.CardsVerListElement{bson.ObjectIdHex(p.SyncTestCardId7), 2},
			// testView.CardsVerListElement{p.SyncTestCardId8, 1},
		}}
}

func (p *publicDataSet) GetDicSearchByTextReqBodyStruct() *testView.ReqDicText {
	p.TextToSearch = "empathy"
	p.SearchSortOption = "createdAt"
	return &testView.ReqDicText{
		WordsText:   p.TextToSearch,
		SortOption:  p.SearchSortOption,
		IsAscending: p.SearchAscending}
}

func (p *publicDataSet) GetDicSearchByIdReqBodyStruct() *testView.ReqDicId {
	return &testView.ReqDicId{
		ParentId:    bson.ObjectIdHex(p.ParentIdInSearch),
		LastId:      bson.ObjectIdHex(p.LastIdInSearch),
		SortOption:  p.SearchSortOption,
		IsAscending: p.SearchAscending}
}

func (p *publicDataSet) GetActivationUrlCodeFromDb(u *userX) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	var ux testView.User
	err = s.DB("mylang").C("users").Find(bson.M{"email": u.Email}).One(&ux)
	defer s.Close()
	p.ActivationUrlCode = ux.ActivationUrlCode
}

func (p *publicDataSet) GetPasswordResettingUrlCodeFromDb(u *userX) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	var ux testView.User
	err = s.DB("mylang").C("users").Find(bson.M{"email": u.Email}).One(&ux)
	defer s.Close()
	p.PasswordResettingUrlCode = ux.PasswordResettingUrlCodes[len(ux.PasswordResettingUrlCodes)-1].PasswordResettingUrlCode
}

func (p *publicDataSet) GetCardVerNoFromDb() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	var c testView.Card
	fmt.Println("p.CardIdOriginal: ", p.CardIdOriginal)
	err = s.DB("mylang").C("cards").Find(bson.M{"_id": bson.ObjectIdHex(p.CardIdOriginal)}).One(&c)
	if err != nil {
		fmt.Println("p.CardIdOriginal card not found: ", err.Error())
	}
	fmt.Println("p.CardIdOriginalVerNo: ", c.VersionNo)
	defer s.Close()
	p.CardIdOriginalVerNo = c.VersionNo
}

func (p *publicDataSet) ClearAllCardsInDb() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	_, _ = s.DB("mylang").C("cards").RemoveAll(bson.M{})
	defer s.Close()
	p.CardIdOriginal = ""
	p.CardIdOriginalVerNo = 0
	p.CardIdDerived = ""
	p.Detail = ""
}

func (p *publicDataSet) PrintAllDicTextInDb() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	s := session.Clone()
	var d []testView.DicText
	err = s.DB("mylang").C("dicTexts").Find(nil).All(&d)
	if err != nil {
		fmt.Println("dicText not got: ", err.Error())
	}
	fmt.Println("dicText in db: ", d)
	defer s.Close()
}

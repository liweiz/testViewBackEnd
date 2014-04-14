package main

import (
	"bytes"
	"fmt"
	//"log"
	"encoding/json"
	//"github.com/go-martini/martini"
	"labix.org/v2/mgo/bson"
	"me/testView/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	TestEmail = "matt.z.lw@gmail.com"
	// password for "matt.z.lw@gmail.com": "aA1~aA1~" or "aA1~aA1~!"
	TestPassword                = "aA1~aA1~"
	TestPasswordV               = "aA1~aA1~!"
	TestUserId                  = "534b2149bd7e5551fc000001"
	TestAccessToken             = "4db0d37878b54392acd054f2ac719a7e"
	TestRefreshToken            = "cf6c657b94f848f5bfa85fd9c7536912"
	TestDeviceUUID              = "zzzzzzzz"
	TestUuidId                  = "534b24b0bd7e555128000001"
	TestRequestId               = "z"
	TestPasswordResettingURLEnd = ""
	TestActivationURLEnd        = "8a643f595f0a48af921df23e267d11bc"
)

// func TestSignUp(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqSignUpOrIn{
// 		TestEmail,
// 		TestPassword,
// 		TestDeviceUUID}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users"
// 	m.Post(url, testView.RequestPreprocessor(testView.SignUp), testView.ProcessedResponseGenerator(testView.SignUp, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestSignIn(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqSignUpOrIn{
// 		"",
// 		"",
// 		TestDeviceUUID}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users/signin"
// 	m.Post(url, testView.GateKeeper(), testView.RequestPreprocessor(testView.SignIn), testView.ProcessedResponseGenerator(testView.SignIn, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.SetBasicAuth(TestEmail, TestPasswordV)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestRenewTokens(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqRenewTokens{
// 		DeviceUUID: TestDeviceUUID,
// 		Tokens:     testView.TokensInCommon{TestAccessToken, TestRefreshToken}}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users/" + TestUserId + "/tokens"
// 	// Use this to simulate use token with wrong userId.
// 	// url := "/users/53498c7abd7e550ce4000002/tokens"
// 	m.Post("/users/:user_id/tokens", testView.GateKeeperExchange(), testView.RequestPreprocessor(testView.RenewTokens), testView.ProcessedResponseGenerator(testView.RenewTokens, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestActivationEmail(t *testing.T) {
// 	m := My()
// 	url := "/users/" + TestUserId + "/activation"
// 	m.Get("/users/:user_id/activation", testView.GateKeeper(), testView.EmailSender(testView.EmailForActivation))
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestClickActivationLink(t *testing.T) {
// 	m := My()
// 	url := "/users/" + TestUserId + "/activation/" + TestActivationURLEnd
// 	m.Get("/users/:user_id/activation/:activation_code", testView.WebPageServer(testView.PageForActivation))
// 	req, _ := http.NewRequest("GET", url, nil)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// This is for reset password when logged in. So gateKeeper needed.
// func TestPasswordResettingEmailByToken(t *testing.T) {
// 	m := My()
// 	url := "/users/" + TestUserId + "/password"
// 	m.Get("/users/:user_id/password", testView.GateKeeper(), testView.EmailSender(testView.EmailForPasswordResetting))
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// This is for forgot-password. No gateKeeper needed, but user has to provide the email to verify if the account exsits. If yes, send the email and let user know. Otherwise, tell the user the account does not exist.
// func TestPasswordResettingEmailByEmail(t *testing.T) {
// 	m := My()
// 	url := "/users/forgotpassword"
// 	reqBodyStruct := testView.ReqForgotPassword{TestEmail}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	body := bytes.NewReader(reqBody)
// 	m.Post("/users/forgotpassword", testView.RequestPreprocessor(testView.ForgotPassword), testView.ProcessedResponseGenerator(testView.ForgotPassword, false))
// 	m.Get("/assets/css/bootstrap.min.css", testView.AssetsServer(testView.BootstrapCssMin))
// 	m.Get("/assets/js/bootstrap.min.js", testView.AssetsServer(testView.BootstrapJsMin))
// 	m.Get("/assets/css/account_activation.css", testView.AssetsServer(testView.CssPageForActivation))
// 	m.Get("/assets/js/account_activation.js", testView.AssetsServer(testView.JsPageForActivation))
// 	m.Get("/assets/css/password_resetting.css", testView.AssetsServer(testView.CssPageForPasswordResetting))
// 	m.Get("/assets/js/password_resetting.js", testView.AssetsServer(testView.JsPageForPasswordResetting))
// 	req, _ := http.NewRequest("POST", url, body)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// Reset password from the uniqueUrl
// func TestChangePassword(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqResetPassword{
// 		TestPasswordV}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users/" + TestUserId + "/password/" + TestPasswordResettingURLEnd
// 	m.Post("/users/:user_id/password/:password_resetting_code", testView.UrlCodeChecker(), testView.RequestPreprocessor(testView.PasswordResetting), testView.ProcessedResponseGenerator(testView.PasswordResetting, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestNewDeviceInfo(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqDeviceInfoNew{
// 		RequestId:  TestRequestId,
// 		DeviceUUID: TestDeviceUUID,
// 		DeviceInfo: testView.DeviceInfoInCommonNew{
// 			BelongTo:   bson.ObjectIdHex("533df9eebd7e554fa0000001"),
// 			DeviceUUID: TestDeviceUUID,
// 			SourceLang: "English",
// 			TargetLang: "简体中文",
// 			SortOption: "timeModifiedDescending",
// 			IsLoggedIn: true,
// 			RememberMe: true}}
// 	fmt.Println(reqBodyStruct)
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users/" + TestUserId + "/deviceinfos"
// 	// Use this to simulate use token with wrong userId.
// 	// url := "/users/53498c7abd7e550ce4000002/deviceinfo"
// 	m.Post("/users/:user_id/deviceinfos", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.NewDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewDeviceInfo, true))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

// func TestUpdateDeviceInfo(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqDeviceInfoNew{
// 		RequestId:  TestRequestId,
// 		DeviceUUID: TestDeviceUUID,
// 		DeviceInfo: testView.DeviceInfoInCommonNew{
// 			BelongTo:   bson.ObjectIdHex(TestUserId),
// 			DeviceUUID: TestDeviceUUID,
// 			SourceLang: "简体中文",
// 			TargetLang: "English",
// 			SortOption: "timeModifiedAscending",
// 			IsLoggedIn: true,
// 			RememberMe: true}}
// 	fmt.Println(reqBodyStruct)
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	// url := "/users/" + TestUserId + "/deviceinfos/" + "TestUuidId"
// 	url := "/users/" + TestUserId + "/deviceinfos/" + TestUuidId
// 	// Use this to simulate use token with wrong userId.
// 	// url := "/users/53498c7abd7e550ce4000002/deviceinfo"
// 	m.Post("/users/:user_id/deviceinfos/:device_id", testView.GateKeeper(), testView.RequestPreprocessor(testView.OneDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneDeviceInfo, true))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

package main

import (
	"bytes"
	"fmt"
	//"log"
	"encoding/json"
	//"github.com/go-martini/martini"
	// "labix.org/v2/mgo/bson"
	"me/testView/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	TestAccessToken  = "3fe35b99cb754bd8b5fe2f1ee203f21e"
	TestRefreshToken = "461cb641b3454a159326a51254869c3a"
	TestDeviceUUID   = "zzzzzzzz"
	TestRequestId    = "zzzzzzzzz"
)

// func TestSignUp(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqSignUpOrIn{
// 		"bbbbbb@gmail.com",
// 		"1111111",
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
// 	m.Post(url, testView.RequestPreprocessor(testView.SignIn), testView.ProcessedResponseGenerator(testView.SignIn, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.SetBasicAuth("bbbbbb@gmail.com", "1111111")
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
// 	url := "/users/533df9eebd7e554fa0000001/deviceinfo"
// 	m.Post("/users/:user_id/deviceinfo", testView.GateKeeper(), testView.RequestPreprocessor(testView.NewDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewDeviceInfo, true))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

func TestRenewTokens(t *testing.T) {
	m := My()
	reqBodyStruct := testView.ReqRenewTokens{
		DeviceUUID: TestDeviceUUID,
		Tokens:     testView.TokensInCommon{TestAccessToken, TestRefreshToken}}
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	os.Stdout.Write(reqBody)
	body := bytes.NewReader(reqBody)
	url := "/users/533df9eebd7e554fa0000001/tokens"
	m.Post("/users/:user_id/tokens", testView.GateKeeperExchange(), testView.RequestPreprocessor(testView.RenewTokens), testView.ProcessedResponseGenerator(testView.RenewTokens, false))
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+TestAccessToken)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	fmt.Println("code:", w.Code)
	os.Stdout.Write(w.Body.Bytes())
}

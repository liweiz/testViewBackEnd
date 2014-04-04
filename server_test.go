package main

import (
	"bytes"
	"fmt"
	//"log"
	"encoding/json"
	//"github.com/go-martini/martini"
	"me/testView/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSignUp(t *testing.T) {
	m := My()
	reqBodyStruct := testView.ReqSignUpOrIn{
		"bbbbbb@gmail.com",
		"1111111",
		"ffffsssaaaaa"}
	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
	os.Stdout.Write(reqBody)
	body := bytes.NewReader(reqBody)
	url := "/users"
	m.Post(url, testView.RequestPreprocessor(testView.SignUp), testView.ProcessedResponseGenerator(testView.SignUp, false))
	req, _ := http.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	fmt.Println("code:", w.Code)
	os.Stdout.Write(w.Body.Bytes())
}

// func TestSignIn(t *testing.T) {
// 	m := My()
// 	reqBodyStruct := testView.ReqSignUpOrIn{
// 		"",
// 		"",
// 		"ffffsssaaaaa"}
// 	reqBody, _ := json.MarshalIndent(reqBodyStruct, "", "	")
// 	os.Stdout.Write(reqBody)
// 	body := bytes.NewReader(reqBody)
// 	url := "/users/signin"
// 	m.Post(url, testView.RequestPreprocessor(testView.SignIn), testView.ProcessedResponseGenerator(testView.SignIn, false))
// 	req, _ := http.NewRequest("POST", url, body)
// 	req.SetBasicAuth("bbbbb@gmail.com", "1111111")
// 	w := httptest.NewRecorder()
// 	m.ServeHTTP(w, req)
// 	fmt.Println("code:", w.Code)
// 	os.Stdout.Write(w.Body.Bytes())
// }

package testView

import (
	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

// func TestPreprocessor(t *testing.T) {
// 	var req *http.Request
// 	var tempParams *martini.Params
// 	Convey("Given a request", t, func() {
// 		reqBodyStruct := ReqSignUpOrIn {
// 			"matt.z.lw@gmail.com",
// 			"1111111",
// 			"ffffsssaaaaa"
// 		}
// 		var reqBody []byte
// 		var err error
// 		reqBody, err = json.Marshal(reqBodyStruct)
// 		req = http.NewRequest("POST", "/users", bytes.NewReader(reqBody))
// 		tempParams = &(martini.Params)(make(map[string]string, 0))

// 		RequestPreprocessor(signUp, req, tempParams, )
// 		})
// }

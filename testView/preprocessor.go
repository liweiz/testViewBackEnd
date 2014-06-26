package testView

import (
	"encoding/json"
	"errors"
	"github.com/go-martini/martini"
	// "io/ioutil"
	"bytes"
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"reflect"
)

// Route categories
const (
	SignUp = iota
	SignIn
	OneUser
	ForgotPassword
	RenewTokens
	NewDeviceInfo
	OneDeviceInfo
	// OneActivationEmail
	OneActivationPage
	// OnePasswordResettingEmail
	OnePasswordResettingPage
	PasswordResetting
	Sync
	NewCard
	OneCard
)

func RequestPreprocessor(route int) martini.Handler {
	return func(req *http.Request, params martini.Params, ctx martini.Context, logger *log.Logger, rw http.ResponseWriter) {
		err := PreprocessRequest(route, req, params, ctx)
		if err != nil {
			HandleReqBodyError(err, logger, rw)
		}
	}
}

// Prepare the incoming and outgoing struct and search criteria for next step
// Both structs of req and res are returned as pointers.
func PreprocessRequest(route int, req *http.Request, params martini.Params, ctx martini.Context) (err error) {
	m := req.Method
	// Get request body and criteria for record(s) searching
	switch route {
	// Sign up
	// case SignUp:
	// Create a new user
	// No sync starts from client after signing up successfully. User choose the lanuage pair and send the newly created deviceInfoInCommon from client.
	// 	if m == "POST" {
	// 		// reqStruct := &ReqSignUpOrIn{}
	// 		resStruct := &ResSignUpOrIn{}
	// 		// err = GetStructFromReq(req, reqStruct)
	// 		if err == nil {
	// 			if len(reqStruct.Password) < 6 {
	// 				err = errors.New("incorrect password format.")
	// 			} else {
	// 				c := bson.M{
	// 					"email": reqStruct.Email}
	// 				PrepareVehicle(ctx, reqStruct, resStruct, c, "", "")
	// 			}
	// 		}
	// 	}
	// case SignIn:
	// 	if m == "POST" {
	// Use martini.Context here as a vehicle to deliver tokens between tokens issued by token server and request_handler.
	// client starts sync right after signing in successfully. If no corresponding deviceInfo on server, get the default deviceInfo with GetDefaultDeviceInfo.
	// 	reqStruct := &ReqSignUpOrIn{}
	// 	resStruct := &ResSignUpOrIn{}
	// 	err = GetStructFromReq(req, reqStruct)
	// 	if err == nil {
	// 		u, err := GetAuthInHeader(req)
	// 		reqStruct.Email = u.Email
	// 		reqStruct.Password = u.Password
	// 		if err == nil {
	// 			c := bson.M{
	// 				"email": reqStruct.Email}
	// 			PrepareVehicle(ctx, reqStruct, resStruct, c, "", "")
	// 		}
	// 	}
	// }
	case ForgotPassword:
		if m == "POST" {
			reqStruct := &ReqForgotPassword{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				c := bson.M{
					"email": reqStruct.Email}
				PrepareVehicle(ctx, reqStruct, nil, c, "", "")
			}
		}
	case RenewTokens:
		if m == "POST" {
			reqStruct := &ReqRenewTokens{}
			resStruct := &ResTokensOnly{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				idToCheck := bson.ObjectIdHex(params["user_id"])
				c := bson.M{
					"belongTo":     idToCheck,
					"accessToken":  reqStruct.Tokens.AccessToken,
					"refreshToken": reqStruct.Tokens.RefreshToken,
					"deviceUUID":   reqStruct.DeviceUUID}
				PrepareVehicle(ctx, reqStruct, resStruct, c, "", "")
			}
		}
	case NewDeviceInfo:
		if m == "POST" {
			reqStruct := &ReqDeviceInfo{}
			resStruct := &ResDeviceInfo{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				PrepareVehicle(ctx, reqStruct, resStruct, nil, reqStruct.RequestId, reqStruct.DeviceUUID)
			}
		}
	case OneDeviceInfo:
		if m == "POST" {
			reqStruct := &ReqDeviceInfo{}
			resStruct := &ResDeviceInfo{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				c := bson.M{
					"_id":        bson.ObjectIdHex(params["device_id"]),
					"belongTo":   bson.ObjectIdHex(params["user_id"]),
					"deviceUUID": reqStruct.DeviceUUID}
				PrepareVehicle(ctx, reqStruct, resStruct, c, reqStruct.RequestId, reqStruct.DeviceUUID)
			}
		}
	case PasswordResetting:
		// Send an email with the password resetting link. E.g., http://www.xxx.com/users/:user_id/passwordresetting/:passwordresetting_code, :passwordresetting_code is used as a unique one time location to reset the password. The location is only valid once. Once being loaded, it becomes invalid. Therefore, to reset password again, user has to use the client to send another email with a new link to resetting password page. The disposable setting of valid link could prevent the page being abused.
		if m == "POST" {
			// This is for resetting password.
			reqStruct := &ReqResetPassword{}
			// Only a 200 header needed for successful request, no need to have body here.
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				idToCheck := bson.ObjectIdHex(params["user_id"])
				c := bson.M{
					"_id": idToCheck}
				PrepareVehicle(ctx, reqStruct, nil, c, "", "")
			}
		}
	case NewCard:
		reqStruct := &ReqCard{}
		resStruct := &ResCards{}
		err = GetStructFromReq(req, reqStruct)
		if err == nil {
			// idToCheck := bson.ObjectIdHex(params["user_id"])
			// c := bson.M{
			// 	"belongTo": idToCheck
			// }
			PrepareVehicle(ctx, reqStruct, resStruct, nil, reqStruct.RequestId, reqStruct.DeviceUUID)
		}
	case OneCard:
		// No need to worry about isDeleted here since the versionNo comparison will take care of that.
		resStruct := &ResCards{}
		idToCheck := bson.ObjectIdHex(params["card_id"])
		belongTo := bson.ObjectIdHex(params["user_id"])
		c := bson.M{
			"_id":       idToCheck,
			"belongTo":  belongTo,
			"isDeleted": false}
		if m == "POST" || m == "DELETE" {
			reqStruct := &ReqCard{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				PrepareVehicle(ctx, reqStruct, resStruct, c, reqStruct.RequestId, reqStruct.DeviceUUID)
			}
			// Add to
		} else if m == "GET" {
			PrepareVehicle(ctx, nil, resStruct, c, "", "")
		}
	case Sync:
		if m == "POST" {
			reqStruct := &ReqSync{}
			resStruct := &ResSync{}
			err = GetStructFromReq(req, reqStruct)
			if err == nil {
				PrepareVehicle(ctx, reqStruct, resStruct, nil, "", reqStruct.DeviceUUID)
			}
		}
	/*
		The request is actually reassemmbled by server to form a query. All the information needed is delivered with the url.
		/dic/:source_lang_code/:target_lang_code/:words_id/:translation_id/:detail_id
	*/

	default:
		err = errors.New("Request not recognized.")
	}
	fmt.Println("request received:")
	var myBuffer bytes.Buffer
	_, _ = myBuffer.ReadFrom(req.Body)
	os.Stdout.Write(myBuffer.Bytes())
	return
}

func GetStructFromReq(req *http.Request, s interface{}) (err error) {
	if req.Body != nil {
		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(s)
		if err != nil {
			fmt.Println("GetStructFromReq err: ", err.Error())
		}
	} else {
		err = errors.New("Request body is nil.")
	}
	return
}

// A struct to store in martini.Context as a vehicle to set/store info from different handlers.
func PrepareVehicle(ctx martini.Context, reqS interface{}, resS interface{}, c bson.M, requestId string, deviceUUID string) {
	v := Vehicle{}

	if reqS != nil {
		v.ReqStruct = reqS
	}
	if resS != nil {
		v.ResStruct = resS
	}
	if c != nil {
		v.Criteria = c
	}
	if requestId != "" {
		v.RequestId = requestId
	}
	if deviceUUID != "" {
		v.DeviceUUID = deviceUUID
	}
	ctx.Map(&v)
}

func GetVehicleContentInContext(ctx martini.Context, fieldName string) reflect.Value {
	v := ctx.Get(reflect.TypeOf(Vehicle{}))
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.FieldByName(fieldName)
}

type Vehicle struct {
	ReqStruct interface{}
	ResStruct interface{}
	Criteria  bson.M
	// Only need to be filled when necessary
	RequestId  string
	DeviceUUID string
}

func AddTextToDic() {

}

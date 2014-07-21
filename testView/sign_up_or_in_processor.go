package testView

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func ProcessedSignUpOrInResponseGenerator(route int) martini.Handler {
	return func(db *mgo.Database, req *http.Request, rw http.ResponseWriter, logger *log.Logger) {
		var result ResSignUpOrIn
		var err error
		if route == SignUp {
			result, err = signUpProcessor(db, logger, req)
		} else if route == SignIn {
			result, err = signInProcessor(db, logger, req)
		} else {
			err = errors.New("Invalid signUp signIn type.")
		}
		var s string
		if err == nil {
			// Send response.
			rw.Header().Set("Content-Type", "application/json")
			var j []byte
			j, err = json.Marshal(result)
			if err == nil {
				// Response size, any usage???
				_, err = rw.Write(j)
				fmt.Println("code:", 200)
				os.Stdout.Write(j)
			}
			if err != nil {
				s = strings.Join([]string{"Failed to generate response, but request has been successfully processed by server.", err.Error()}, "=> ")
			}
		}
		if err != nil {
			if s == "" {
				WriteLog(err.Error(), logger)
				fmt.Println("err to log: ", err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			} else {
				WriteLog(s, logger)
				http.Error(rw, s, http.StatusServiceUnavailable)
			}
		}
	}
}

// Only requests that pass the gateKeeper are processed here. This indicates a not found err here means user not activated.
// This is the revised one used to process user with Lang pair
func signUpProcessor(db *mgo.Database, logger *log.Logger, r *http.Request) (result ResSignUpOrIn, err error) {
	a, err := GetAuthInHeader(r)
	if err == nil {
		if len(a.Password) < 6 || len(a.Password) > 20 {
			err = errors.New("incorrect password format.")
		} else {
			var aUser User
			// Check if already in use
			err = db.C("users").Find(bson.M{
				"email": a.Email}).One(&aUser)
			if err == mgo.ErrNotFound {
				fmt.Println("existing user not found")
				// Reset err = nil
				err = nil
				// Create a new user
				newId := bson.NewObjectId()
				uuid.SwitchFormat(uuid.Clean, false)
				uniqueUrlCode := uuid.NewV4().String()
				var hashedPassword []byte
				hashedPassword, err = bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
				if err == nil {
					reqStruct := &ReqSignUpOrIn{}
					err = GetStructFromReq(r, reqStruct)
					if err == nil {
						docToSave := bson.M{
							// UserInCommon
							"activated":  false,
							"email":      a.Email,
							"_id":        newId,
							"versionNo":  1,
							"sourceLang": reqStruct.SourceLang,
							"targetLang": reqStruct.TargetLang,

							//Non UserInCommon part
							"lastModified":      time.Now().UnixNano(),
							"createdAt":         time.Now().UnixNano(),
							"isDeleted":         false,
							"password":          hashedPassword,
							"activationUrlCode": uniqueUrlCode}
						err = db.C("users").Insert(docToSave)
						if err != nil {
							fmt.Println("bson.M to insert err: ", err.Error())
						}
						if err == nil {
							// Get and put the new user to response body.
							var rr UserInCommon
							err = db.C("users").Find(bson.M{
								"_id": newId}).Select(GetSelector(SelectUserInCommon)).One(&rr)
							if err == nil {
								// Send activation email after user is successfully created on server.
								uniqueUrl := path.Join(r.URL.String(), newId.Hex(), "activation", uniqueUrlCode)
								// Ignore potential error for sending activation email here. In case it fails and user does not received the mail, user can trigger sending the mail later.
								_ = GenerateEmail(EmailForActivation, uniqueUrl, a.Email)
								result.User = rr
								// fake a structFromReq with real email and password: &ReqSignUpOrIn{}
								s := &ReqSignUpOrIn{}
								s.Email = a.Email
								s.Password = a.Password
								s.DeviceUUID = r.Header.Get("X-REMOLET-DEVICE-ID")
								var r1 TokensInCommon
								r1, err = SetGetDeviceTokens(rr.Id, s, db)
								if err == nil {
									result.Tokens = r1
								}
							}
						}
					}
				}
			} else if err == nil {
				err = errors.New("User already exists.")
			}
		}
	}
	return
}

// This is the original one, which does not process user with lang pair
func signUpProcessorOriginal(db *mgo.Database, logger *log.Logger, r *http.Request) (result ResSignUpOrIn, err error) {
	a, err := GetAuthInHeader(r)
	if err == nil {
		if len(a.Password) < 6 || len(a.Password) > 20 {
			err = errors.New("incorrect password format.")
		} else {
			var aUser User
			// Check if already in use
			err = db.C("users").Find(bson.M{
				"email": a.Email}).One(&aUser)
			if err == mgo.ErrNotFound {
				fmt.Println("existing user not found")
				// Reset err = nil
				err = nil
				// Create a new user
				newId := bson.NewObjectId()
				uuid.SwitchFormat(uuid.Clean, false)
				uniqueUrlCode := uuid.NewV4().String()
				var hashedPassword []byte
				hashedPassword, err = bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
				if err == nil {
					docToSave := bson.M{
						// UserInCommon
						"activated": false,
						"email":     a.Email,
						"_id":       newId,
						"versionNo": 1,

						//Non UserInCommon part
						"lastModified":      time.Now().UnixNano(),
						"createdAt":         time.Now().UnixNano(),
						"isDeleted":         false,
						"password":          hashedPassword,
						"activationUrlCode": uniqueUrlCode}
					err = db.C("users").Insert(docToSave)
					if err != nil {
						fmt.Println("bson.M to insert err: ", err.Error())
					}
					if err == nil {
						// Get and put the new user to response body.
						var rr UserInCommon
						err = db.C("users").Find(bson.M{
							"_id": newId}).Select(GetSelector(SelectUserInCommon)).One(&rr)
						if err == nil {
							// Send activation email after user is successfully created on server.
							uniqueUrl := path.Join(r.URL.String(), newId.Hex(), "activation", uniqueUrlCode)
							// Ignore potential error for sending activation email here. In case it fails and user does not received the mail, user can trigger sending the mail later.
							_ = GenerateEmail(EmailForActivation, uniqueUrl, a.Email)
							result.User = rr
							// fake a structFromReq: &ReqSignUpOrIn{}
							s := &ReqSignUpOrIn{}
							s.Email = a.Email
							s.Password = a.Password
							s.DeviceUUID = r.Header.Get("X-REMOLET-DEVICE-ID")
							var r1 TokensInCommon
							r1, err = SetGetDeviceTokens(rr.Id, s, db)
							if err == nil {
								result.Tokens = r1
							}
						}
					}
				}
			} else if err == nil {
				err = errors.New("User already exists.")
			}
		}
	}
	return
}

func signInProcessor(db *mgo.Database, logger *log.Logger, r *http.Request) (result ResSignUpOrIn, err error) {
	a, err := GetAuthInHeader(r)
	if err == nil {
		var rr UserInCommon
		err = db.C("users").Find(bson.M{
			"email": a.Email}).Select(GetSelector(SelectUserInCommon)).One(&rr)
		if err == nil {
			result.User = rr
			// fake a structFromReq: &ReqSignUpOrIn{}
			s := &ReqSignUpOrIn{}
			s.Email = a.Email
			s.Password = a.Password
			s.DeviceUUID = r.Header.Get("X-REMOLET-DEVICE-ID")
			var r1 TokensInCommon
			r1, err = SetGetDeviceTokens(rr.Id, s, db)
			if err == nil {
				result.Tokens = r1
			}
		}
	}
	return
}

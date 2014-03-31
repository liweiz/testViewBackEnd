package testView

import (
	"encoding/base64"
	"errors"
	"github.com/codegangsta/martini"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

func PrimaryAuthHandler(db *mgo.Database, r *http.Request, rw martini.ResponseWriter, logger *log.Logger) {
	isValid, err := ValidatePrimaryAuth(db, r)
	if !isValid {
		if err != nil {
			WriteLog(err.Error(), logger)
			if err.Error() == "Incorrect password" || err.Error() == "Token expired" {
				http.Error(rw, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == "Invalid authorization header" {
				http.Error(rw, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func ValidatePrimaryAuth(db *mgo.Database, r *http.Request) (isValid bool, err error) {
	var auth *AuthInHeader
	auth, err = GetAuthInHeader(r)
	if err == nil {
		var isMatched bool
		isMatched, err = MatchPrimaryAuth(auth, db)
		if isMatched {
			isValid = true
		}
	}
	return
}

// PrimaryAuth means password or accessToken. isMatched simply means it is the same as the one stored in db. But it may already be expired.
func MatchPrimaryAuth(auth *AuthInHeader, db *mgo.Database) (isMatched bool, err error) {
	if auth.AuthType == "Basic" {
		var user User
		err = db.C("users").Find(bson.M{"email": auth.Email}).One(&user)
		if err == nil {
			if auth.Password == user.Password {
				isMatched = true
			} else {
				err = errors.New("Incorrect password")
			}
		}
	}
	if auth.AuthType == "Bearer" {
		var myDeviceTokens DeviceTokens
		err = db.C("deviceTokens").Find(bson.M{"accessToken": auth.AccessToken}).One(&myDeviceTokens)
		// AccessToken exists, check expiration next.
		// AccessToken expired, check JSON body for refresh token, ask for refreshToken. This occurs at URL: /users/:user_id/deviceinfo/:device_id/token, case renewTokens. "Access token expired. Refresh token needed"
		if err == nil {
			if CheckTokenExpiration(&myDeviceTokens) {
				err = errors.New("Token expired")
			} else {
				isMatched = true
			}
		}
	}
	return
}

func CheckTokenExpiration(d *DeviceTokens) (isExpired bool) {
	if d.AccessTokenExpireAt <= time.Now().UnixNano() {
		// AccessToken expired
		isExpired = true
	}
	return
}

// Return authorization header data
func GetAuthInHeader(r *http.Request) (a *AuthInHeader, err error) {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		s := strings.SplitN(auth, " ", 2)
		if s[0] == "Basic" {
			var b []byte
			b, err = base64.StdEncoding.DecodeString(s[1])
			if err == nil {
				pair := strings.SplitN(string(b), ":", 2)
				if len(pair) == 2 {
					a = &AuthInHeader{AuthType: s[0], Email: pair[0], Password: pair[1]}
					return
				}
			}
		} else if s[0] == "Bearer" {
			a = &AuthInHeader{AuthType: s[0], AccessToken: s[1]}
			return
		}
	}
	if err == nil {
		err = errors.New("Invalid authorization header")
	}
	return
}

// Parse basic authentication header
type AuthInHeader struct {
	AuthType    string
	Email       string
	Password    string
	AccessToken string
}

// AccessTokenGen generates access tokens
type TokensGen interface {
	GenerateTokens(generateRefresh bool) (accessToken string, refreshToken string, err error)
}

/*
GenerateAccessToken generates version 4 UUID access and refresh tokens
No base64 encode needed for tokens. See https://github.com/bshaffer/oauth2-server-php/issues/100:
On casual reading of "The OAuth 2.0 Authorization Protocol: Bearer Tokens"* I've encountered several people (including myself) who have made the assumption that the name b64token implies that some kind of base64 encoding/decoding on the access token is taking place between the client and RS. Digging a bit deeper in to "HTTP/1.1, part 7: Authentication"**, however, I see that b64token is just an ABNF syntax definition allowing for characters typically used in base64, base64url, etc.. So the b64token doesn't define any encoding or decoding but rather just defines what characters can be used in the part of the Authorization header that will contain the access token.
But email and password sent from clients are still base64 encoded to avoid clear text.
*/
func GenerateTokens(generateRefresh bool) (accessToken string, refreshToken string) {
	uuid.SwitchFormat(uuid.Clean, false)
	accessToken = uuid.NewV4().String()
	if generateRefresh {
		refreshToken = uuid.NewV4().String()
	}
	return
}

package tokenServer

import (
        "net/http"
        "time"
        "strings"
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "github.com/codegangsta/martini"
        "testView"
        "log"
        "errors"
)

func PrimaryAuthHandler(db *mgo.Database, r *http.Request, rw *martini.ResponseWriter, logger *log.Logger) {
    isValid, err := ValidatePrimaryAuth(db, r)
    if !isValid {
        if err != nil {
            WriteLog(err.Error(), logger)
            if err.Error() == "Incorrect password" || err.Error() == "Token expired" {
                http.Error(rw, err.Error(), StatusUnauthorized)
            } else if err.Error() == "Invalid authorization header" {
                http.Error(rw, err.Error(), StatusBadRequest)
            } else {
                http.Error(rw, err.Error(), StatusInternalServerError)
            }
        }
    }
}

func ValidatePrimaryAuth(db *mgo.Database, r *http.Request) (isValid bool, err error) { 
    var auth *AuthInHeader
    auth, err := GetAuthInHeader(r)
    if !err {
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
        err = db.C("users").Find(bson.M{"email": auth.Email}).One(&user)
        if !err {
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
        if !err {
            if CheckTokenExpiration(myDeviceTokens) {
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
            if !err {
                pair := strings.SplitN(string(b), ":", 2)
                if len(pair) == 2 {
                    a = &AuthInHeader{AuthType: s[0], Email: pair[0], Password: pair[1]}
                    return
                }
            }
        } else if s[0] == "Bearer" {
            a = &AuthInHeader{AuthType: s[0], Token: s[1]}
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
    AuthType string
    Email string
    Password string
    Token string
}

// AccessTokenGen generates access tokens
type TokensGen interface {
    GenerateTokens(generateRefresh bool) (accessToken string, refreshToken string, err error)
}

// GenerateAccessToken generates base64-encoded UUID access and refresh tokens
func GenerateTokens(generateRefresh bool) (accessToken string, refreshToken string) {
    accessToken = uuid.New()
    accessToken = base64.StdEncoding.EncodeToString([]byte(accessToken))
    if generateRefresh {
        refreshToken = uuid.New()
        refreshToken = base64.StdEncoding.EncodeToString([]byte(refreshToken))
    }
    return
}
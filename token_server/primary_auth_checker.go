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

func ValidatePrimaryAuth(db *mgo.Database, r *http.Request) (err error, isValid bool) { 
    auth, err1 := GetAuthInHeader(r)
    err = err1
    if err == nil {
        err2, isMatched := MatchPrimaryAuth(auth, db)
        err = err2
        if isMatched {
            isValid = true
        }
    }
    return
}

// PrimaryAuth means password or accessToken. isMatched simply means it is the same as the one stored in db. But it may already be expired.
func MatchPrimaryAuth(auth *AuthInHeader, db *mgo.Database) (err error, isMatched bool) {
    if auth.AuthType == "Basic" {
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
        if len(s) != 2 {
            err = errors.New("Invalid authorization header")
        } else if s[0] == "Basic" {
            b, err := base64.StdEncoding.DecodeString(s[1])
            if err != nil {
                return nil, err
            }
            pair := strings.SplitN(string(b), ":", 2)
            if len(pair) != 2 {
                return nil, errors.New("Invalid authorization message")
            }
            return &AuthInHeader{AuthType: s[0], Email: pair[0], Password: pair[1]}, nil
        } else if s[0] == "Bearer" {
            return &AuthInHeader{AuthType: s[0], Token: s[1]}, nil
        } else {
            return nil, errors.New("Invalid authorization header")
        }
    } else {
        err = errors.New("Invalid authorization header")
    }
    return
}
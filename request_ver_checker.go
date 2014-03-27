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

func CheckReqVerNo(db *mgo.Database, route int, criteria bson.M, structFromReq interface{}, req *http.Request, structForRes interface{}, params *martini.Params)
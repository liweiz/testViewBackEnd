package testView

import (
		"net/http"
		"encode/json"
		"github.com/codegangsta/martini"
		"labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
        "net/url"
        "error"
        "reflect"
)

func HandleReqBodyError(err error, logger *log.Logger, rw *martini.ResponseWriter) {
    WriteLog(err.Error(), logger)
    // Either no body found or internal error
    if err.Error() == "Request body is nil." {
        http.Error(rw, err.Error(), StatusBadRequest)
    } else {
        http.Error(rw, err.Error(), StatusInternalServerError)
    }
}

// Response for unsuccessfully proccessed request. No body in this case.
func WriteUnauthorizedRes(errorMsg string, errorDetail string, w *martini.ResponseWriter) {
    // StatusUnauthorized = 401
    writer.WriteHeader(writer.StatusUnauthorized)
    // Something like this:
    // WWW-Authenticate: Bearer error="invalid_token",
    //                   error_description="The access token expired"
    s := []string{"Bearer error=\"", errorMsg, "\", error_description=\"", errorDetail, "\""}
    writer.Header().Set("WWW-Authenticate", strings.Join(s, ""))
}

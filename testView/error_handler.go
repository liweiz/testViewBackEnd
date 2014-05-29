package testView

import (
	//"github.com/go-martini/martini"
	"log"
	"net/http"
	"strings"
)

func HandleReqBodyError(err error, logger *log.Logger, rw http.ResponseWriter) {
	WriteLog(err.Error(), logger)
	// Either no body found or internal error
	if err.Error() == "Request body is nil." {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// Response for unsuccessfully proccessed request. No body in this case.
func WriteUnauthorizedRes(errorMsg string, errorDetail string, w http.ResponseWriter) {
	// StatusUnauthorized = 401
	w.WriteHeader(http.StatusUnauthorized)
	// Something like this:
	// WWW-Authenticate: Bearer error="invalid_token",
	//                   error_description="The access token expired"
	s := []string{"Bearer error=\"", errorMsg, "\", error_description=\"", errorDetail, "\""}
	w.Header().Set("WWW-Authenticate", strings.Join(s, ""))
}

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
// Get request body and criteria for record(s) searching
// Make decision
// Generate response based on decision made

/*
0. Token check. Put it here since tokens are stored in header. No need to unmarshal body.
1. Preprocess the request to get the request body into struct, prepare the criteria for searching and struct for response use.
0. In the case of login/refreshToken, issue a new set of tokens/ask for relogin. Leave it here since email/password/refreshToken are all stored in body. Hence, unmarshal needed.
2. Use the search criteria to operate in db and update the resStruct in martini.Context
*/

func ProcessRequest(route int, req *http.Request, db *mgo.Database, rw martini.ResponseWriter) {
	reqStruct, resStruct, criteria, err := PreProcessRequest(route, req, params *martini.Params)
	var decision int
	if err != nil {
		// Send response for err
		decision = InternalServerError
	} else {
		decision = MakeDecision(db, route, criteria, reqStruct, req, resStruct, params)
	}
	
}
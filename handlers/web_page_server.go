package testView

import (
	"errors"
	"github.com/go-martini/martini"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"path"
	// "text/template"
	"time"
)

// webPageCategory
const (
	PageForActivation = iota
	PageForPasswordResetting
)

const MyHostAddress string = "localhost:3000"

func WebPageServer(pCategory int) martini.Handler {
	return func(db *mgo.Database, rw http.ResponseWriter, req *http.Request, params martini.Params, logger *log.Logger) {
		err := ServeWebPage(db, pCategory, rw, req, params)
		if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

func ServeWebPage(db *mgo.Database, pCategory int, rw http.ResponseWriter, req *http.Request, params martini.Params) (err error) {
	var pageFileName string
	var pageFilePath string
	n := CssAndJsNeeded{}
	n.BaseCss = path.Join(MyHostAddress, "pages/css/bootstrap.min.css")
	n.BaseJs = path.Join(MyHostAddress, "pages/css/bootstrap.min.css")
	switch pCategory {
	case PageForActivation:
		var aUser User
		err = db.C("users").Find(bson.M{"_id": bson.ObjectIdHex(params["user_id"])}).One(&aUser)
		if err == nil {
			if !aUser.Activated {
				err = db.C("users").Update(bson.M{"_id": bson.ObjectIdHex(params["user_id"])}, bson.M{
					"$set": bson.M{"activated": true, "lastModified": time.Now().UnixNano()},
					"$inc": bson.M{"versionNo": 1}})
			}
		}
		if err != nil {
			return
		}
		pageFileName = "account_activation.html"
		pageFilePath = path.Join("pages/account_activation", pageFileName)
		n.Css = path.Join(MyHostAddress, "assets/css/account_activation.css")
		n.Js = path.Join(MyHostAddress, "assets/js/account_activation.js")
	case PageForPasswordResetting:
		pageFileName = "password_resetting.html"
		pageFilePath = path.Join("pages/password_resetting", pageFileName)
		n.Css = path.Join(MyHostAddress, "assets/css/password_resetting.css")
		n.Js = path.Join(MyHostAddress, "assets/js/password_resetting.js")
	default:
		err = errors.New("No page file set for this purpose.")
		return
	}
	var t *template.Template
	t, err = template.ParseFiles("pages/layout/webpage_layout.html", pageFilePath)
	if err == nil {
		err = t.Execute(rw, n)
	}
	return
}

type CssAndJsNeeded struct {
	BaseCss string
	BaseJs  string
	Css     string
	Js      string
}

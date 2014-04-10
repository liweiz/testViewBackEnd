package testView

import (
	"github.com/go-martini/martini"
	"html/template"
	"log"
	"net/http"
	"path"
)

// webPageCategory
const (
	PageForActivation = iota
	PageForPasswordResetting
)

func WebPageServer(pCategory int) martini.Handler {
	return func(rw http.ResponseWriter, req *http.Request, params martini.Params, logger *log.Logger) {
		err := ServeWebPage(pCategory, rw, req, params)
		if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

func ServeWebPage(pCategory int, rw http.ResponseWriter, req *http.Request, params martini.Params) (err error) {
	var pageFileName string
	var pageFilePath string
	switch pageCategory {
	case PageForActivation:
		pageFileName = "account_activation.html"
		pageFilePath = path.Join("pages/account_activation", pageFileName)
	case PageForPasswordResetting:
		pageFileName = "password_resetting.html"
		pageFilePath = path.Join("pages/password_resetting", pageFileName)
	default:
		err = errors.New("No page file set for this purpose.")
		return
	}
	var t *template.Template
	t, err = template.ParseFiles("pages/layout.html", pageFilePath)
	if err == nil {
		err = t.Execute(rw, nil)
	}
	return
}

package testView

import (
	"errors"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

// assetCode
const (
	BootstrapCssMin = iota
	BootstrapJsMin
	CssPageForActivation
	JsPageForActivation
	CssPageForPasswordResetting
	JsPageForPasswordResetting
)

func AssetsServer(assetCode int) martini.Handler {
	return func(rw http.ResponseWriter, req *http.Request, logger *log.Logger) {
		assetNeeded, err := ServeAsset(assetCode)
		fmt.Println(assetNeeded)
		if err != nil {
			WriteLog(err.Error(), logger)
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		} else {
			http.ServeFile(rw, req, assetNeeded)
		}
	}
}

func ServeAsset(assetCode int) (AssetNeeded string, err error) {
	switch assetCode {
	case BootstrapCssMin:
		AssetNeeded = "/pages/css/bootstrap.min.css"
	case BootstrapJsMin:
		AssetNeeded = "/pages/js/bootstrap.min.js"
	case CssPageForActivation:
		AssetNeeded = "/pages/account_activation/account_activation.css"
	case JsPageForActivation:
		AssetNeeded = "/pages/account_activation/account_activation.js"
	case CssPageForPasswordResetting:
		AssetNeeded = "/pages/password_resetting/password_resetting.css"
	case JsPageForPasswordResetting:
		AssetNeeded = "/pages/password_resetting/password_resetting.js"
	default:
		err = errors.New("No asset for this purpose.")
	}
	return
}

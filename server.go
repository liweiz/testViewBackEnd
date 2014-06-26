package main

import (
	//"fmt"
	"github.com/go-martini/martini"
	//"io"
	"labix.org/v2/mgo"
	"me/testViewPKG/testView"
	//"net/http"
	//"os"
)

func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("mylang"))
		defer s.Close()
		c.Next()
	}
}

type MyMartini struct {
	*martini.Martini
	martini.Router
}

func My() *MyMartini {
	m := martini.New()
	r := martini.NewRouter()
	// Setup middleware
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(DB())
	m.Use(r.Handle)
	m.MapTo(r, (*martini.Routes)(nil))
	return &MyMartini{m, r}
}

func main() {
	// go func() {
	// 	if err:= http.ListenAndServe(":8000", http.Handlerfunc(func(w http.ResponseWriter, r *http.Request) {
	// 		http.Error(w, "https scheme is required", http.StatusBadRequest)
	// 		})); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// if err := http.ListenAndServeTLS(":8001", "cert.pem", "key.pem", m); err != nil {
	// 	log.Fatal(err)
	// }

	m := My()
	// Exchange for a new set of tokens. renewTokens
	m.Post("/users/:user_id/tokens", testView.GateKeeperExchange(), testView.RequestPreprocessor(testView.RenewTokens), testView.ProcessedResponseGenerator(testView.RenewTokens, false))
	// Update a user's device settings, e.g., language pair/sort option. OneDeviceInfoLang
	m.Post("/users/:user_id/deviceinfos/:device_id", testView.GateKeeper(), testView.RequestPreprocessor(testView.OneDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneDeviceInfo, true))
	// Create a new deviceInfo
	m.Post("/users/:user_id/deviceinfos", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.NewDeviceInfo), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewDeviceInfo, true))
	// Sync cards and user. sync. No request id needed for sync. Everytime a sync request received by server, server responses and client take the result to decide what to do next.
	m.Post("/users/:user_id/sync", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.Sync), testView.ProcessedResponseGenerator(testView.Sync, false))
	// Update a new card. oneCard
	m.Post("/users/:user_id/cards/:card_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	// Create a new card. newCard
	m.Post("/users/:user_id/cards", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.NewCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewCard, true))
	// Change a user's password.
	m.Post("/users/:user_id/password/:password_resetting_code", testView.UrlCodeChecker(), testView.RequestPreprocessor(testView.PasswordResetting), testView.ProcessedResponseGenerator(testView.PasswordResetting, false))
	// Send an email with url to change a user's password in the case of forgot-password.
	m.Post("/users/forgotpassword", testView.RequestPreprocessor(testView.ForgotPassword), testView.ProcessedResponseGenerator(testView.ForgotPassword, false))
	// User signs in. signIn
	m.Post("/users/signin", testView.GateKeeper(), testView.ProcessedSignUpOrInResponseGenerator(testView.SignIn))
	// Sign up a new user. signUp
	m.Post("/users", testView.ProcessedSignUpOrInResponseGenerator(testView.SignUp))

	// Get a card. oneCard
	//m.Get("/users/:user_id/cards/:card_id")
	// Activate a user. activation
	m.Get("/users/:user_id/activation/:activation_code", testView.WebPageServer(testView.PageForActivation))
	// Send an email with url to activate a user.
	m.Get("/users/:user_id/activation", testView.GateKeeper(), testView.EmailSender(testView.EmailForActivation))
	// Serve the webpage to reset a user's password. passwordResetting
	m.Get("/users/:user_id/password/:password_resetting_code", testView.WebPageServer(testView.PageForPasswordResetting))
	// Send an email with url to change a user's password.
	m.Get("/users/:user_id/password", testView.GateKeeper(), testView.EmailSender(testView.EmailForPasswordResetting))
	// Get user
	m.Get("/users/:user_id", testView.GateKeeper(), testView.ProcessedResponseGenerator(testView.OneUser, false))

	// Delete a card. oneCard
	m.Delete("/users/:user_id/cards/:card_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	// Get all cards
	// m.Get("/users/:user_id/cards")

	m.Post("/dic/:source_lang_code/:target_lang_code/text/:user_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.DicTextSearcher())
	m.Post("/dic/:source_lang_code/:target_lang_code/id/:user_id", testView.GateKeeper(), testView.NonActivationBlocker(), testView.DicIdSearcher())

	// Serve assets
	m.Get("/assets/css/bootstrap.min.css", testView.AssetsServer(testView.BootstrapCssMin))
	m.Get("/assets/js/bootstrap.min.js", testView.AssetsServer(testView.BootstrapJsMin))
	m.Get("/assets/css/account_activation.css", testView.AssetsServer(testView.CssPageForActivation))
	m.Get("/assets/js/account_activation.js", testView.AssetsServer(testView.JsPageForActivation))
	m.Get("/assets/css/password_resetting.css", testView.AssetsServer(testView.CssPageForPasswordResetting))
	m.Get("/assets/js/password_resetting.js", testView.AssetsServer(testView.JsPageForPasswordResetting))
	m.Run()
}

// The regex to check for the requested format (allows an optional trailing
// slash).
// var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
// func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
// 	// Get the format extension
// 	matches := rxExt.FindStringSubmatch(r.URL.Path)
// 	ft := ".json"
// 	if len(matches) > 1 {
// 		// Rewrite the URL without the format extension
// 		l := len(r.URL.Path) - len(matches[1])
// 		if strings.HasSuffix(r.URL.Path, "/") {
// 			l--
// 		}
// 		r.URL.Path = r.URL.Path[:l]
// 		ft = matches[1]
// 	}
// 	// Inject the requested encoder
// 	switch ft {
// 	case ".xml":
// 		c.MapTo(xmlEncoder{}, (*Encoder)(nil))
// 		w.Header().Set("Content-Type", "application/xml")
// 	case ".text":
// 		c.MapTo(textEncoder{}, (*Encoder)(nil))
// 		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	default:
// 		c.MapTo(jsonEncoder{}, (*Encoder)(nil))
// 		w.Header().Set("Content-Type", "application/json")
// 	}
// }

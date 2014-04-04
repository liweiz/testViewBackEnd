package main

import (
	//"fmt"
	"github.com/go-martini/martini"
	//"io"
	"labix.org/v2/mgo"
	"me/testView/handlers"
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

	// Exchange for a new set of tokens. renewTokens
	m := My()
	m.Post("/users/:user_id/deviceinfo/:device_id/token", testView.RequestPreprocessor(testView.RenewTokens), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.RenewTokens, false))
	// Update a user's device settings, e.g., language pair/sort option. OneDeviceInfoSortOption
	m.Post("/users/:user_id/deviceinfo/:device_id/SortOption", testView.RequestPreprocessor(testView.OneDeviceInfoSortOption), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneDeviceInfoSortOption, false))
	// Update a user's device settings, e.g., language pair/sort option. OneDeviceInfoLang
	m.Post("/users/:user_id/deviceinfo/:device_id/Lang", testView.RequestPreprocessor(testView.OneDeviceInfoLang), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneDeviceInfoLang, false))
	// Sync cards and user. sync
	m.Post("/users/:user_id/sync", testView.RequestPreprocessor(testView.Sync), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.Sync, false))
	// Update a new card. oneCard
	m.Post("/users/:user_id/cards/:card_id", testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	// Create a new card. newCard
	m.Post("/users/:user_id/cards", testView.RequestPreprocessor(testView.NewCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.NewCard, true))
	// Change a user's email. oneUser
	m.Post("/users/:user_id", testView.RequestPreprocessor(testView.OneUser), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneUser, false))
	// User signs in. signIn
	m.Post("/users/signin", testView.RequestPreprocessor(testView.SignIn), testView.ProcessedResponseGenerator(testView.SignIn, false))
	// Sign up a new user. signUp
	m.Post("/users", testView.RequestPreprocessor(testView.SignUp), testView.ProcessedResponseGenerator(testView.SignUp, false))

	// Get a card. oneCard
	//m.Get("/users/:user_id/cards/:card_id")
	// Activate a user. activation
	//m.Get("/users/:user_id/activation/:activation_code")
	// Reset a user's password. passwordResetting
	//m.Get("/users/:user_id/passwordresetting/:passwordresetting_code")

	// Delete a card. oneCard
	m.Delete("/users/:user_id/cards/:card_id", testView.RequestPreprocessor(testView.OneCard), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.OneCard, true))
	// Get all cards
	// m.Get("/users/:user_id/cards")

	// Get context list based on words, translation and detail. dicDetail
	//m.Get("/dic/:sourcelang/:targetlang/:words_id/:translation_id/:detail_id")
	// Get detail list based on words and translation. dicTranslation
	//m.Get("/dic/:sourcelang/:targetlang/:words_id/:translation_id")
	// Get translation list based on words. dicWords
	m.Post("/dic/:sourcelang/:targetlang", testView.RequestPreprocessor(testView.DicWords), testView.ReqIdChecker(), testView.ProcessedResponseGenerator(testView.DicWords, false))

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

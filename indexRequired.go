package testView

import {
	"log"
	"net/http"
	"regexp"
	"strings"
	"database/sql"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/auth"
	"labix.org/v2/mgo"
}

func GetDicIndexReady(db *mgo.Database) {
	db.C("dictionary").
}

type DicIndex mgo.Index {
	Key:        []string{"sourceLang", "targetLang", "context", "target", "translation", "detail", "contextLength", "targetLength", "translationLength", "detailLength"},
	Unique:     false,
	DropDups:   false,
	Background: true,
	Sparse:     false,
}

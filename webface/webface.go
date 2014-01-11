/*
Package webface provides web interface to view collected data.
It will not work unless you copy static files and templates into
directory with an executable:
 - webface/templates/*
 - webface/static/*
Web server must serve static files!
*/
package webface

import (
	"bitbucket.org/kardianos/osext"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/tsukanov/steamhistory/storage"
	"github.com/tsukanov/steamhistory/storage/analysis"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"
)

// Start starts FastCGI server at 127.0.0.1:9000
func Start() {
	log.Println("Starting server...")
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal("Failed to start server!", err)
	}
	fcgi.Serve(l, makeRouter())
}

// StartDev starts a development server at localhost:8080
func StartDev() {
	log.Println("Starting development server (localhost:8080)...")
	http.ListenAndServe(":8080", makeRouter())
}

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/history/{appid:[0-9]+}.json", historyHandler)
	r.HandleFunc("/top/", topHandler)
	r.HandleFunc("/top/daily.json", dailyTopHandler)
	r.HandleFunc("/search", searchHandler)
	r.HandleFunc("/about/", aboutHandler)
	return r
}

var mc *memcache.Client = memcache.New("localhost:11211")

/*
 * Handlers
 */

// basicHandler just loads specified template, combines it with base template
// and writes result into ResponseWriter.
func basicHandler(w http.ResponseWriter, r *http.Request, file string) {
	exeloc, err := osext.ExecutableFolder()
	t, err := template.ParseFiles(
		exeloc+"webface/templates/base.html",
		exeloc+"webface/templates/"+file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	basicHandler(w, r, "home.html")
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	key := "history_" + strconv.Itoa(appId)
	it, err := mc.Get(key)
	var b []byte
	if err == nil {
		b = it.Value
	} else {
		name, err := storage.GetName(appId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		history, err := storage.AllUsageHistory(appId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		type jason struct {
			Name    string     `json:"name"`
			History [][2]int64 `json:"history"`
		}
		result := jason{
			Name:    name,
			History: history,
		}
		b, err = json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 1800}) // 1800 sec = 30 min
		if err != nil {
			log.Println(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	basicHandler(w, r, "top.html")
}

func dailyTopHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := analysis.MostPopularAppsToday()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	key := "top"
	it, err := mc.Get(key)
	var b []byte
	if err == nil {
		b = it.Value
	} else {
		b, err = json.Marshal(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 1800}) // 1800 sec = 30 min
		if err != nil {
			log.Println(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	query, ok := queries["q"]
	if !ok {
		http.Error(w, "No query", http.StatusBadRequest)
		return
	}

	h := md5.New()
	key := fmt.Sprintf("%x", h.Sum([]byte(query[0])))
	it, err := mc.Get(key)
	var b []byte
	if err == nil {
		b = it.Value
	} else {
		results, err := storage.Search(query[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		b, err = json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 43200}) // 43200 sec = 12 hours
		if err != nil {
			log.Println(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	basicHandler(w, r, "about.html")
}

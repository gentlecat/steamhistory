// Package webface provides web interface to view collected data.
package webface

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/steamhistory/core/analysis"
	"github.com/steamhistory/core/apps"
	"github.com/steamhistory/core/usage"
)

// Start starts FastCGI server at 127.0.0.1:9000
func Start() {
	log.Println("Starting server...")
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal("Failed to start server!", err)
	}
	log.Println("Listening on 127.0.0.1:9000...")
	fcgi.Serve(l, makeRouter())
}

// StartDev starts development server at localhost:8080
func StartDev() {
	log.Println("Starting development server (localhost:8080)...")
	http.ListenAndServe(":8080", makeRouter())
}

func makeRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/apps", appsHandler)
	r.HandleFunc("/apps/popular", dailyPopularHandler)
	r.HandleFunc("/history/{appid:[0-9]+}", historyHandler)
	return r
}

var mc *memcache.Client = memcache.New("localhost:11211")

/*
 * Handlers
 */

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("See API documentation.\n"))
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	key := "history_" + strconv.Itoa(appId)
	it, err := mc.Get(key)
	var b []byte
	if err == nil {
		b = it.Value
	} else {
		name, err := apps.GetName(appId)
		if err != nil {
			http.Error(w, "Internal error.", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		history, err := usage.AllUsageHistory(appId)
		if err != nil {
			http.Error(w, "Internal error.", http.StatusInternalServerError)
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
			http.Error(w, "Internal error.", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 1800}) // 1800 sec = 30 min
		if err != nil {
			log.Println(err)
		}
	}

	queries := r.URL.Query()
	callback, ok := queries["callback"]
	if ok {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s(%s)", callback[0], b)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func dailyPopularHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := analysis.MostPopularAppsToday()
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
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
			http.Error(w, "Internal error.", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 1800}) // 1800 sec = 30 min
		if err != nil {
			log.Println(err)
		}
	}

	queries := r.URL.Query()
	callback, ok := queries["callback"]
	if ok {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s(%s)", callback[0], b)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func appsHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	query, ok := queries["q"]
	if !ok {
		// TODO: Return all apps
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
		results, err := apps.Search(query[0])
		if err != nil {
			http.Error(w, "Internal error.", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		b, err = json.Marshal(results)
		if err != nil {
			http.Error(w, "Internal error.", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = mc.Set(&memcache.Item{Key: key, Value: b, Expiration: 43200}) // 43200 sec = 12 hours
		if err != nil {
			log.Println(err)
		}
	}

	callback, ok := queries["callback"]
	if ok {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s(%s)", callback[0], b)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

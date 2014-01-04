package webface

import (
	"bitbucket.org/kardianos/osext"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/tsukanov/steaminfo-go/storage"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"
)

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/history/{appid:[0-9]+}.json", historyHandler)
	r.HandleFunc("/search", searchHandler)
	r.HandleFunc("/about/", aboutHandler)
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	exeloc, err := osext.ExecutableFolder()
	t, err := template.ParseFiles(
		exeloc+"webface/templates/base.html",
		exeloc+"webface/templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = t.Execute(w, r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	mc := memcache.New("localhost:11211")
	it, err := mc.Get("history_" + string(appId))
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
		mc.Set(&memcache.Item{Key: "history_" + string(appId), Value: b, Expiration: 1800}) // 1800 sec = 30 min
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
	results, err := storage.Search(query[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	b, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	exeloc, err := osext.ExecutableFolder()
	t, err := template.ParseFiles(
		exeloc+"webface/templates/base.html",
		exeloc+"webface/templates/about.html")
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

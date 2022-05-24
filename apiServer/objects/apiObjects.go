package objects

import (
	"dsxy/es"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func del(w http.ResponseWriter, r *http.Request) {
	fmt.Println("del debug")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

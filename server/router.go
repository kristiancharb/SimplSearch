package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type IndexReq struct {
	Name string
}

type DocumentReq struct {
	Contents string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to SimplSearch!")
}

func newIndexHandler(w http.ResponseWriter, r *http.Request) {
	var body IndexReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Created new index named: %v", body.Name)
}

func newDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var body DocumentReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexName := mux.Vars(r)["name"]
	fmt.Fprintf(w,
		"Created new document for index %v with contents:\n%v",
		indexName,
		body.Contents)
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/index", newIndexHandler).Methods("POST")
	router.HandleFunc("/index/{name}", newDocumentHandler).Methods("POST")
	return router
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kristiancharb/SimplSearch/search"
	"net/http"
)

type IndexWrapper struct {
	store *search.IndexStore
}

type IndexReq struct {
	Name string
}

type DocumentReq struct {
	Title    string
	Contents string
}

type SearchReq struct {
	Query string
	Start int
	End   int
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	indexStore := search.NewIndexStore()
	index := IndexWrapper{indexStore}

	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/index", index.newIndexHandler).Methods("POST")
	router.HandleFunc("/index/{name}", index.newDocumentHandler).Methods("POST")
	router.HandleFunc("/search/{name}", index.searchHandler).Methods("POST")
	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to SimplSearch!!!")
}

func (index *IndexWrapper) newIndexHandler(w http.ResponseWriter, r *http.Request) {
	var body IndexReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index.store.NewIndex(body.Name)
	fmt.Fprintf(w, "Created new index named: %v", body.Name)
}

func (index *IndexWrapper) newDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var body DocumentReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexName := mux.Vars(r)["name"]
	index.store.AddDocument(indexName, body.Title, body.Contents, -1)
	index.store.UpdateIndex(indexName)
	fmt.Fprintf(w,
		"Created new document for index %v",
		indexName,
	)
}

func (index *IndexWrapper) searchHandler(w http.ResponseWriter, r *http.Request) {
	var body SearchReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if body.Start < 0 || body.End < 0 || body.Start > body.End {
		http.Error(w, "Invalid range", http.StatusBadRequest)
		return
	}
	indexName := mux.Vars(r)["name"]
	queryResponse := index.store.Search(indexName, body.Query, body.Start, body.End)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queryResponse)
}

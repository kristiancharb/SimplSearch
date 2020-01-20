package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kristiancharb/SimplSearch/search"
	"net/http"
)

type IndexReq struct {
	Name string
}

type DocumentReq struct {
	Contents string
}

type SearchReq struct {
	Query string
}

type IndexStore struct {
	store map[string]*search.Index
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	indexStore := IndexStore{make(map[string]*search.Index)}

	router.HandleFunc("/", handler).Methods("GET")
	router.HandleFunc("/index", indexStore.newIndexHandler).Methods("POST")
	router.HandleFunc("/index/{name}", indexStore.newDocumentHandler).Methods("POST")
	router.HandleFunc("/search/{name}", indexStore.searchHandler).Methods("POST")
	return router
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to SimplSearch!!!")
}

func (indexStore *IndexStore) newIndexHandler(w http.ResponseWriter, r *http.Request) {
	var body IndexReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexStore.store[body.Name] = &search.Index{
		Name:  body.Name,
		Terms: make(map[string]*search.TermInfo),
	}
	fmt.Fprintf(w, "Created new index named: %v", body.Name)
}

func (indexStore *IndexStore) newDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var body DocumentReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexName := mux.Vars(r)["name"]
	indexStore.store[indexName].AddDocument(body.Contents)
	fmt.Fprintf(w,
		"Created new document for index %v",
		indexName,
	)
}

func (indexStore *IndexStore) searchHandler(w http.ResponseWriter, r *http.Request) {
	var body SearchReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	indexName := mux.Vars(r)["name"]
	indexStore.store[indexName].Search(body.Query)
}

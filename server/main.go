package main

import (
	// "fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handler)
	http.ListenAndServe(":8000", router)
}

package main

import (
	"net/http"
)

func main() {
	router := newRouter()
	http.ListenAndServe(":8000", router)
}

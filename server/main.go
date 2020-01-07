package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8000"
	router := newRouter()
	fmt.Printf("Server started on port %v\n", port)
	http.ListenAndServe(":8000", router)
}

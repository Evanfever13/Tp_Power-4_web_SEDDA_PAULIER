package main

import "net/http"

func mainRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/error", ErrorHandler)
}

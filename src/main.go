package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	LoadTemplates()

	mux := http.NewServeMux()

	mainRouter(mux)

	addr := "localhost:8080"
	fmt.Printf("Serve on : http://%s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}

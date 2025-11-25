package main

import (
	"fmt"
	"log"
	"net/http"
)

// === Fonction Main ===
func main() {
	// Chargement des templates et fichiers statiques
	LoadTemplates()

	// Initialisation du routeur
	mux := http.NewServeMux()
	mainRouter(mux)

	// Démarrage du serveur
	addr := "localhost:8080"
	fmt.Printf("Serve on : http://%s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}

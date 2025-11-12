package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {

	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "Home", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/gameinit", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "GameInit", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	fichierserveur := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fichierserveur))

	fmt.Print("Le Serveur est Lancé sur : http://localhost:8000/")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}

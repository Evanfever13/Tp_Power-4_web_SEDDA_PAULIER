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

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "home", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	fileServer := http.FileServer(http.Dir("./../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}

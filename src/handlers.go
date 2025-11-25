package main

import (
	"fmt"
	"net/http"
)

// === Handlers ===

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "Home", nil)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("code"))
	fmt.Println(r.FormValue("message"))
	fmt.Fprintf(w, "Une erreur est survenue lors du traitement de votre requete.")
}

/*
func handlers() {

	// Routes

	http.HandleFunc("/gameinit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Impossible de lire le formulaire", http.StatusBadRequest)
				return
			}
			var NewGame = Game{
				Player1Name: r.FormValue("joueur1"),
				Player2Name: r.FormValue("joueur2"),
				Gameboard:   [6][7]string{},
			}
			NewGame.InitPlayer()
			fmt.Println("Noms:", NewGame.Player1Name, "vs", NewGame.Player2Name)
			fmt.Println("Plateau initialisé")
			http.Redirect(w, r, "/gameplay", http.StatusSeeOther)
			return
		}
		RenderTemplate(w, r, "GameInit", nil)
	})

	http.HandleFunc("/game/play", func(w http.ResponseWriter, r *http.Request) {
		// Utilise la logique de jeu
		GamePlay(w, r, temp, &NewGame)
	})

	http.HandleFunc("/ScoreBoard", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Affichage du ScoreBoard")

		ScoreBoard = sortScoreBoard(ScoreBoard)
		if err := temp.ExecuteTemplate(w, "ScoreBoard", ScoreBoard); err != nil {
			fmt.Printf("erreur exec : %v\n", err.Error())
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/gameend", func(w http.ResponseWriter, r *http.Request) {
		// Utilise la logique de jeu
		GamePlay(w, r, temp, &NewGame)
	})
	// Fichiers statiques
	fichierserveur := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fichierserveur))

	// Serveur
	fmt.Println("Le Serveur est Lancé sur : http://localhost:8000/")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}
*/

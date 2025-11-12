package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func WinCheck(board [6][7]string) bool {
	// Vérification horizontale
	for r := 0; r < 6; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r][c+1] && board[r][c] == board[r][c+2] && board[r][c] == board[r][c+3] {
				return true
			}
		}
	}
	// Vérification verticale
	for c := 0; c < 7; c++ {
		for r := 0; r < 3; r++ {
			if board[r][c] != "" && board[r][c] == board[r+1][c] && board[r][c] == board[r+2][c] && board[r][c] == board[r+3][c] {
				return true
			}
		}
	}
	// Vérification diagonale (de haut-gauche à bas-droite)
	for r := 0; r < 3; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r+1][c+1] && board[r][c] == board[r+2][c+2] && board[r][c] == board[r+3][c+3] {
				return true
			}
		}
	}
	// Vérification diagonale (de bas-gauche à haut-droite)
	for r := 3; r < 6; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r-1][c+1] && board[r][c] == board[r-2][c+2] && board[r][c] == board[r-3][c+3] {
				return true
			}
		}
	}
	return false
}
func main() {
	// Initialisation du jeu
	type Game struct {
		Player1    string
		Player2    string
		Gameboard  [6][7]string
		Playerturn string
	}
	NewGame := Game{}

	// Chargement des templates HTML
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	// Gestion des routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "Home", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/gameinit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Impossible de lire le formulaire", http.StatusBadRequest)
				return
			}
			joueur1 := r.FormValue("joueur1")
			joueur2 := r.FormValue("joueur2")
			NewGame = Game{
				Player1:    joueur1,
				Player2:    joueur2,
				Gameboard:  [6][7]string{},
				Playerturn: joueur1,
			}
			fmt.Println("Noms des joueurs :", NewGame.Player1, "et", NewGame.Player2)
			fmt.Println("Initialisation du plateau de jeu...", NewGame.Gameboard)
			fmt.Println("C'est au tour de", NewGame.Player1)
			fmt.Println("Démarrage d'une nouvelle partie...")

			http.Redirect(w, r, "/gameplay", http.StatusSeeOther)
			return
		}

		if err := temp.ExecuteTemplate(w, "GameInit", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/gameplay", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "gameplay", NewGame); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	// Gestion des fichiers statiques
	fichierserveur := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fichierserveur))

	// Lancement du serveur
	fmt.Print("Le Serveur est Lancé sur : http://localhost:8000/")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println("Erreur serveur:", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// === Structure du jeu ===
type Game struct {
	PlayerWinner string
	Player1Name  string
	Player2Name  string
	Player1Sym   string
	Player2Sym   string
	Gameboard    [6][7]string
	Playerturn   string // contient le symbole courant ("X" ou "O")
}
type ScoreEntry struct {
	Name  string
	Score int
}

var ScoreBoard []ScoreEntry

// === Mise à jour du ScoreBoard ===
func updateScoreBoard(winner string) {
	for i, entry := range ScoreBoard {
		if entry.Name == winner {
			ScoreBoard[i].Score++
			return
		}
	}
	ScoreBoard = append(ScoreBoard, ScoreEntry{Name: winner, Score: 1})
}

// === Fonction de tri du ScoreBoard ===
func sortScoreBoard(sb []ScoreEntry) []ScoreEntry {
	// Retourne une slice triée par score décroissant.
	entries := make([]ScoreEntry, 0, len(sb))
	for _, entry := range sb {
		entries = append(entries, entry)
	}

	// Tri decroissant par score, puis par nom
	for i := 1; i < len(entries); i++ {
		key := entries[i]
		j := i - 1
		for j >= 0 && (entries[j].Score < key.Score || (entries[j].Score == key.Score && entries[j].Name > key.Name)) {
			entries[j+1] = entries[j]
			j--
		}
		entries[j+1] = key
	}
	return entries
}

// === Vérification victoire ===
func WinCheck(board [6][7]string) bool {
	// Horizontale
	for r := 0; r < 6; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r][c+1] && board[r][c] == board[r][c+2] && board[r][c] == board[r][c+3] {
				return true
			}
		}
	}
	// Verticale
	for c := 0; c < 7; c++ {
		for r := 0; r < 3; r++ {
			if board[r][c] != "" && board[r][c] == board[r+1][c] && board[r][c] == board[r+2][c] && board[r][c] == board[r+3][c] {
				return true
			}
		}
	}
	// Diagonale (de haut-gauche à bas-droite)
	for r := 0; r < 3; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r+1][c+1] && board[r][c] == board[r+2][c+2] && board[r][c] == board[r+3][c+3] {
				return true
			}
		}
	}
	// Diagonale (de bas-gauche à haut-droite)
	for r := 3; r < 6; r++ {
		for c := 0; c < 4; c++ {
			if board[r][c] != "" && board[r][c] == board[r-1][c+1] && board[r][c] == board[r-2][c+2] && board[r][c] == board[r-3][c+3] {
				return true
			}
		}
	}
	return false
}

// === Initialisation des symboles et du tour ===
func (g *Game) InitPlayer() {

	g.Player1Sym, g.Player2Sym = "X", "O"
	g.Playerturn = g.Player1Sym
	fmt.Println("Le joueur 1 commence avec X")

}

// === Changement de joueur ===
func (g *Game) switchPlayer() {
	if g.Playerturn == g.Player1Sym {
		g.Playerturn = g.Player2Sym
		fmt.Println("Tour du joueur 2 (", g.Player2Sym, ")")
	} else {
		g.Playerturn = g.Player1Sym
		fmt.Println("Tour du joueur 1 (", g.Player1Sym, ")")
	}
}

// === Ajout d’un jeton ===
func (g *Game) AddJeton(colonne int) bool {
	if colonne < 0 || colonne >= 7 {
		return false
	}
	for i := 5; i >= 0; i-- {
		if g.Gameboard[i][colonne] == "" {
			g.Gameboard[i][colonne] = g.Playerturn
			return true
		}
	}
	return false
}

// === Gameplay ===
func GamePlay(w http.ResponseWriter, r *http.Request, temp *template.Template, g *Game) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Impossible de lire le formulaire", http.StatusBadRequest)
			return
		}
		colStr := r.FormValue("colonne")
		if colStr == "" {
			colStr = r.FormValue("Play")
		}

		var col int
		_, err := fmt.Sscanf(colStr, "%d", &col)
		if err != nil {
			http.Error(w, "Colonne invalide", http.StatusBadRequest)
			return
		}
		col--

		if !g.AddJeton(col) {
			fmt.Println("Colonne pleine ou invalide")
		} else {
			if WinCheck(g.Gameboard) {
				fmt.Println("Victoire du symbole", g.Playerturn) // Test + Vérif du Winner
				if g.Playerturn == g.Player1Sym {
					g.PlayerWinner = g.Player1Name
				} else {
					g.PlayerWinner = g.Player2Name
				}
				updateScoreBoard(g.PlayerWinner)
				if err := temp.ExecuteTemplate(w, "GameEnd", g); err != nil {
					fmt.Println("Template Victory manquant, retour au gameplay") // En cas de manque d'info/Templates
				} else {
					return
				}
			} else {
				g.switchPlayer()
			}
		}
	}

	if err := temp.ExecuteTemplate(w, "gameplay", g); err != nil {
		http.Error(w, "Erreur Templates", http.StatusInternalServerError) // En cas de manque d'info/Templates
	}
}

// === Main ===
func main() {
	rand.Seed(time.Now().UnixNano())

	NewGame := Game{}

	// Chargement des templates HTML
	temp, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	// Routes
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
			NewGame = Game{
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
		if err := temp.ExecuteTemplate(w, "GameInit", nil); err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/gameplay", func(w http.ResponseWriter, r *http.Request) {
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

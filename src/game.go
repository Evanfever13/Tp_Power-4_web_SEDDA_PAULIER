package main

import (
	"fmt"
	"log"
	"net/http"
)

// === Structures ===
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

// === Variables Globales ===
var NewGame Game
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
	entries = append(entries, sb...)
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

// === Vérification si le plateau est plein ===
func isBoardFull(board [6][7]string) bool {
	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			if board[r][c] == "" {
				return false
			}
		}
	}
	return true
}

// === Initialisation des symboles et du tour ===
func (g *Game) InitPlayer() {
	fmt.Println(g)
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
func GamePlay(w http.ResponseWriter, r *http.Request, g *Game) {
	// Parse le formulaire
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	// Récupère la colonne jouée
	colStr := r.FormValue("colonne")
	if colStr == "" {
		colStr = r.FormValue("Play")
	}

	// Convertit en entier
	var col int
	_, err := fmt.Sscanf(colStr, "%d", &col)
	if err != nil {
		log.Println("Colonne invalide")
	}
	col--

	// Ajoute le jeton
	if !g.AddJeton(col) {
		log.Println("Colonne pleine ou invalide")
	} else {

		// Vérifie la victoire
		if WinCheck(g.Gameboard) {
			log.Println("Victoire du symbole", g.Playerturn)
			if g.Playerturn == g.Player1Sym {
				g.PlayerWinner = g.Player1Name
			} else {
				g.PlayerWinner = g.Player2Name
			}
			updateScoreBoard(g.PlayerWinner)
			http.Redirect(w, r, "/game/end", http.StatusSeeOther)

			// Vérifie si le plateau est plein (match nul)
		} else if isBoardFull(g.Gameboard) {
			log.Println("Match nul")
			g.PlayerWinner = "Nul"
			http.Redirect(w, r, "/game/end", http.StatusSeeOther)
		} else {
			// Change de joueur
			fmt.Println(g)
			g.switchPlayer()
		}
	}
}

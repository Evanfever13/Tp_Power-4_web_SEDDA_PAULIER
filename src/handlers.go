package main

import (
	"fmt"
	"net/http"
)

// === HomeHandler ===
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "Home", nil)
}

// === ErrorHandler ===
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	type ErrorData struct {
		Message string
		Code    string
	}

	data := ErrorData{
		Message: r.FormValue("message"),
		Code:    r.FormValue("code"),
	}
	RenderTemplate(w, r, "Error", data)
}

// === GameInitHandler ===
func GameInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
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
		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
		return
	}
	RenderTemplate(w, r, "GameInit", nil)
}

// === GamePlayHandler ===
func GamePlayHandler(w http.ResponseWriter, r *http.Request) {
	GamePlay(w, r, &NewGame)
	fmt.Println(&NewGame)
	RenderTemplate(w, r, "GamePlay", &NewGame)
}

// === GameEndHandler ===
func GameEndHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "GameEnd", &NewGame)
}

// === ScoreBoardHandler ===
func ScoreBoardHandler(w http.ResponseWriter, r *http.Request) {
	ScoreBoard = sortScoreBoard(ScoreBoard)
	RenderTemplate(w, r, "ScoreBoard", ScoreBoard)
}

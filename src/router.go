package main

import "net/http"

// === Main Router ===
func mainRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/error", ErrorHandler)
	mux.HandleFunc("/game/init", GameInitHandler)
	mux.HandleFunc("/game/play", GamePlayHandler)
	mux.HandleFunc("/game/end", GameEndHandler)
	mux.HandleFunc("/scoreboard", ScoreBoardHandler)

	// === Initialisation des fichiers statiques ===
	fichierserveur := http.FileServer(http.Dir("./../assets"))
	mux.Handle("/static/", http.StripPrefix("/static/", fichierserveur))

}

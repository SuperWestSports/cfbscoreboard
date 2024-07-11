package main

import "net/http"

func handleGames(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := gamesTpl.Execute(w, ByConference())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleFeatured(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := featuredTpl.Execute(w, ByFeatured())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleStandings(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := standingsTpl.Execute(w, ByConferenceStandings())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := indexTpl.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
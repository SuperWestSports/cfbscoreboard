package main

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"
)

var (
	gameData     = make(map[int]Game)
	recordData   = make(map[int]Record)
	dataMutex    sync.RWMutex
	gamesTpl     *template.Template
	featuredTpl  *template.Template
	standingsTpl *template.Template
	prod         = false
)

func main() {
	if prod {
		requestGamesUrl := "https://api.collegefootballdata.com/scoreboard"
		requestRecordsUrl := "https://api.collegefootballdata.com/records"
		yearRecord := "?year=2023"
		fetchGamesData(requestGamesUrl)
		fetchRecordsData(requestRecordsUrl, yearRecord)

	} else {
		fetchSampleGamesData("livegamedata.json")
		fetchSampleRecordsData("samplerecords.json")
	}

	var err error
	gamesTpl, err = template.ParseFiles("templates/games.html")
	templateError(err)

	featuredTpl, err = template.ParseFiles("templates/featured.html")
	templateError(err)

	standingsTpl, err = template.ParseFiles("templates/standings.html")
	templateError(err)

	gameData = formatDate(gameData)
	http.HandleFunc("/games", handleGames)
	http.HandleFunc("/featured", handleFeatured)
	http.HandleFunc("/standings", handleStandings)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

func templateError(err error) {
	if err != nil {
		fmt.Println("Error loading template:", err)
	}
}

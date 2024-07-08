package main

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"
)

var (
	requestUrl = "https://api.collegefootballdata.com/scoreboard"
	games []Game
	records []Record
	gameData  = make(map[int]Game)
	recordData = make(map[int]Record)
	dataMutex   sync.RWMutex
	gamesTpl    *template.Template
	featuredTpl *template.Template
	standingsTpl *template.Template
	prod = false
)

func main() {
	if prod {
		req := request()
		res := response(req)
		go fetch(res)
	} else {
			gamesJSON := readSample("livegamedata.json")
			unmarshalSampleGames(gamesJSON)
			loadSampleGames(games)
			
			recordsJSON := readSample("samplerecords.json")
			unmarshalSampleRecords(recordsJSON)
			loadSampleRecords(records)

	}
	
	var err error
	gamesTpl, err = template.ParseFiles("games.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return
	}

	featuredTpl, err = template.ParseFiles("featured.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return
	}

	standingsTpl, err = template.ParseFiles("standings.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return
	}

	gameData = formatDate(gameData)
	http.HandleFunc("/games", handleGames)
	http.HandleFunc("/featured", handleFeatured)
	http.HandleFunc("/standings", handleStandings)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

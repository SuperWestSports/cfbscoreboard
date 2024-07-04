package main

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"
)

var (
	requestUrl = "https://api.collegefootballdata.com/scoreboard"
	gameData  = make(map[int]Game)
	dataMutex   sync.RWMutex
	gamesTpl    *template.Template
	featuredTpl *template.Template
	prod = false
)

func main() {
	if prod {
		req := request()
		res := response(req)
		go fetch(res)
	} else {
			loadSample()
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

	gameData = formatDate(gameData)
	http.HandleFunc("/games", handleGames)
	http.HandleFunc("/featured", handleFeatured)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

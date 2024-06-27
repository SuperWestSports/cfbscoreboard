package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	//req := request()
	//res := response(req)
	//go fetch(res)
	loadSample()

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

	formatDate(gameData)
	http.HandleFunc("/games", handleGames)
	http.HandleFunc("/featured", handleFeatured)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

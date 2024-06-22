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
	tpl, err = template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return
	}

	http.HandleFunc("/games", handleGames)
	http.ListenAndServe(":8080", nil)
}


// fetch data should only be used in real-time for the current week of games
// otherwise scheduled can be requested and saved. Final can be requested and saved
// I need to build this logic with date/time and a few functions

// Aug 24 - Dec 15

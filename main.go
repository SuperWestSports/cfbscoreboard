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
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

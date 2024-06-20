package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	prod := false
	if prod {
		go fetchData()
	} else {
		
		err := loadSampleData()
		if err != nil {
			fmt.Println("Error loading sample data:", err)
			return
		}
	}

	var err error
		tpl, err = template.ParseFiles("template.html")
		if err != nil {
			fmt.Println("Error loading template:", err)
			return
		}

	http.HandleFunc("/games", handleGames)
	log.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

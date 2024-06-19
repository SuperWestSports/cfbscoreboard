package main

import (
	"html/template"
	"net/http"
)
func createTemplate() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := WeekPageData{
			PageTitle: "College Football Scoreboard",
			Weeks: []Week{
				{Title: "Week 0", WeekNumber: 0},
				{Title: "Week 1", WeekNumber: 1},
				{Title: "Week 2", WeekNumber: 2},
			},
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8080", nil)
}

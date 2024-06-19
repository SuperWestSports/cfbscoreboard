package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Week struct {
	Title      string
	WeekNumber int
	Scoreboard []Game
}

type WeekPageData struct {
	PageTitle string
	Weeks     []Week
}

type Game struct {
	id                 int
	season             int
	week               int
	season_type        string
	start_date         string
	start_time_tbd     bool
	completed          bool
	neutral_site       bool
	conference_game    bool
	attendance         int
	venue_id           int
	venue              string
	home_id            int
	home_team          string
	home_conference    string
	home_division      string
	home_points        int
	home_line_scores   []int
	home_post_win_prob float64
	home_pregame_elo   int
	home_postgame_elo  int
	away_id            int
	away_team          string
	away_conference    string
	away_division      string
	away_points        int
	away_line_scores   []int
	away_post_win_prob float64
	away_pregame_elo   int
	away_postgame_elo  int
	excitement_index   float64
}

func main() {
	section := "games"
	year := "2023"
	week := "1"
	KEY := "RMy62JITIczdOcIcgpVpLhfOsl4BlOFvLWsW/NGM/ZgiCcbL3bRK7JnbISToCImy"
	URL := fmt.Sprintf("https://api.collegefootballdata.com/%s?year=%s&week=%s", section, year, week)
	req, err := http.NewRequest(
		http.MethodGet,
		URL,
		nil,
	)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+KEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading HTTP resonse body: %v", err)
	}
	err = os.WriteFile("data/2023week1.json", responseBytes, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//log.Println("We go the response:", string(responseBytes))

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

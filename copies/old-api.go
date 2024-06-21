package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func apiSettings() string {
	section := "games"
	year := "2023"
	week := "1"
	return fmt.Sprintf("https://api.collegefootballdata.com/%s?year=%s&week=%s", section, year, week)
}

func formatRequest() *http.Request {
	KEY := "RMy62JITIczdOcIcgpVpLhfOsl4BlOFvLWsW/NGM/ZgiCcbL3bRK7JnbISToCImy"
	URL := apiSettings()
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
	return req
}

func sendRequest(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}
	if err != nil {
		log.Fatalf("error reading HTTP resonse body: %v", err)
	}
	return resp
}

type Game struct {
	ID             int    `json:"id"`
	Season         int    `json:"season"`
	Week           int    `json:"week"`
	SeasonType     string `json:"season_type"`
	StartDate      string `json:"start_date"`
	StartTimeTBD   bool   `json:"start_time_tbd"`
	Completed      bool   `json:"completed"`
	NeutralSite    bool   `json:"neutral_site"`
	ConferenceGame bool   `json:"conference_game"`
	//Attendance        int       `json:"attendance"`
	VenueID        int    `json:"venue_id"`
	Venue          string `json:"venue"`
	HomeID         int    `json:"home_id"`
	HomeTeam       string `json:"home_team"`
	HomeConference string `json:"home_conference"`
	HomeDivision   string `json:"home_division"`
	HomePoints     int    `json:"home_points"`
	HomeLineScores []int  `json:"home_line_scores"`
	//HomePostWinProb   float64   `json:"home_post_win_prob"`
	//HomePregameElo    int       `json:"home_pregame_elo"`
	//HomePostgameElo   int       `json:"home_postgame_elo"`
	AwayID         int    `json:"away_id"`
	AwayTeam       string `json:"away_team"`
	AwayConference string `json:"away_conference"`
	AwayDivision   string `json:"away_division"`
	AwayPoints     int    `json:"away_points"`
	AwayLineScores []int  `json:"away_line_scores"`
	//AwayPostWinProb   float64   `json:"away_post_win_prob"`
	//AwayPregameElo    int       `json:"away_pregame_elo"`
	//AwayPostgameElo   int       `json:"away_postgame_elo"`
	//ExcitementIndex   float64   `json:"excitement_index"`
	//Highlights        string    `json:"highlights"`
	//Notes             string    `json:"notes"`
}

var (
	gameData []Game
	//gameData  = make(map[int]Game)
	dataMutex sync.RWMutex
	tpl       *template.Template
)

func fetchData() {
	for {
		req := formatRequest()
		resp := sendRequest(req)

		var games []Game
		err := json.NewDecoder(resp.Body).Decode(&games)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		dataMutex.Lock()
		for _, game := range games {
			gameData[game.ID] = game
		}
		dataMutex.Unlock()

		time.Sleep(5 * time.Minute) // Adjust the interval as needed
	}
}

func loadSampleData() error {
	data, err := os.ReadFile("data/2023week1.json")
	if err != nil {
		return err
	}

	var games []Game
	err = json.Unmarshal(data, &games)
	if err != nil {
		return err
	}

	dataMutex.Lock()
	gameData = games
	dataMutex.Unlock()

	return nil
}

func handleGames(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	w.Header().Set("Content-Type", "text/html")
	err := tpl.Execute(w, gameData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

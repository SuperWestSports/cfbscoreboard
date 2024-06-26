package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

type Venue struct {
	Name  string `json:"name"`
	City  string `json:"city"`
	State string `json:"state"`
}

type Team struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Conference     string `json:"conference"`
	Classification string `json:"classification"`
	Points         int    `json:"points"`
}

type Weather struct {
	Temperature   string `json:"temperature"`
	Description   string `json:"description"`
	WindSpeed     string `json:"windSpeed"`
	WindDirection string `json:"windDirection"`
}

type Betting struct {
	Spread        string `json:"spread"`
	OverUnder     string `json:"overUnder"`
	HomeMoneyline int    `json:"homeMoneyline"`
	AwayMoneyline int    `json:"awayMoneyline"`
}

type Game struct {
	ID             int     `json:"id"`
	StartDate      string  `json:"startDate"`
	StartTimeTBD   bool    `json:"startTimeTBD"`
	TV             string  `json:"tv"`
	NeutralSite    bool    `json:"neutralSite"`
	ConferenceGame bool    `json:"conferenceGame"`
	Status         string  `json:"status"`
	Period         int     `json:"period"`
	Clock          string  `json:"clock"`
	Situation      string  `json:"situation"`
	Possession     string  `json:"possession"`
	Venue          Venue   `json:"venue"`
	HomeTeam       Team    `json:"homeTeam"`
	AwayTeam       Team    `json:"awayTeam"`
	Weather        Weather `json:"weather"`
	Betting        Betting `json:"betting"`
}

var (
	requestUrl = "https://api.collegefootballdata.com/scoreboard"
	gameData   []Game
	//gameData  = make(map[int]Game)
	dataMutex sync.RWMutex
	tpl       *template.Template
)

func request() *http.Request {
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CFB_KEY"))
	return req
}

func response(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func fetch(resp *http.Response) {
	for {
		var games []Game
		err := json.NewDecoder(resp.Body).Decode(&games)
		if err != nil {
			fmt.Println("Error decoding JSON", err)
			time.Sleep(10 * time.Minute)
			continue
		}

		dataMutex.Lock()
		for _, game := range games {
			gameData[game.ID] = game
		}
		dataMutex.Unlock()
		time.Sleep((5 * time.Minute))
	}
}

func loadSample() error {
	data, err := os.ReadFile("data/livegamedata.json")
	if err != nil {
		return err
	}

	var games []Game
	err = json.Unmarshal(data, &games)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dataMutex.Lock()
	gameData = games
	dataMutex.Unlock()

	return nil
}

func ByConference(gameData []Game) map[string][]Game {
	sortedConf := make(map[string][]Game)
	conferences := []string{"ACC", "American Athletic", "Big 12", "Big Ten",
		"Conference USA", "FBS Independent", "Mid-American",
		"Mountain West", "Pac-12", "SEC", "Sun Belt"}
	for i := range gameData {
		for j := range conferences {
			if gameData[i].AwayTeam.Conference == conferences[j] || gameData[i].HomeTeam.Conference == conferences[j] {
				sortedConf[conferences[j]] = append(sortedConf[conferences[j]], gameData[i])
			}
		}
	}
	return sortedConf
}

func handleGames(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	w.Header().Set("Content-Type", "text/html")
	err := tpl.Execute(w, ByConference(gameData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

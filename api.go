package main

import (
	"encoding/json"
	"fmt"

	//"io"
	"log"
	"net/http"
	"os"

	//"path/filepath"
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
	Description   string  `json:"description"`
	WindSpeed     string `json:"windSpeed"`
	WindDirection string `json:"windDirection"`
}

type Betting struct {
	Spread        string `json:"spread"`
	OverUnder     string `json:"overUnder"`
	HomeMoneyline int     `json:"homeMoneyline"`
	AwayMoneyline int     `json:"awayMoneyline"`
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
	//gameData = make(map[int]Game)
	dataMutex sync.RWMutex
	tpl       *template.Template
)

func request() *http.Request {
	KEY := "RMy62JITIczdOcIcgpVpLhfOsl4BlOFvLWsW/NGM/ZgiCcbL3bRK7JnbISToCImy"
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+KEY)
	return req
}

func response(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//writeSample(resp)
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

// func writeSample(resp *http.Response) {
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}
// 	err = os.MkdirAll(filepath.Dir("data/scoreboard-preseason.json"), os.ModePerm)
// 	if err != nil {
// 		fmt.Println("Error creating directory:", err)
// 		return
// 	}

// 	err = os.WriteFile("data/scoreboard-preseason.json", body, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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

func handleGames(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	w.Header().Set("Content-Type", "text/html")
	err := tpl.Execute(w, gameData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

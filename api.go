package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func fetchGamesData(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CFB_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var games []Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	dataMutex.Lock()
	for _, game := range games {
		gameData[game.ID] = game
	}
	dataMutex.Unlock()

	return nil
}

func fetchRecordsData(url, year string) error {
	req, err := http.NewRequest("GET", url+year, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CFB_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var records []Record
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	dataMutex.Lock()
	for _, record := range records {
		recordData[record.ID] = record
	}
	dataMutex.Unlock()

	return nil
}

func fetchSampleGamesData(filename string) error {
	data, err := os.ReadFile("data/" + filename)
	if err != nil {
		fmt.Println(err)
	}

	var games []Game
	err = json.Unmarshal(data, &games)
	if err != nil {
		log.Fatal(err)
	}

	dataMutex.Lock()
	for _, game := range games {
		gameData[game.ID] = game
	}
	dataMutex.Unlock()
	return nil
}

func fetchSampleRecordsData(filename string) error {
	data, err := os.ReadFile("data/" + filename)
	if err != nil {
		fmt.Println(err)
	}

	var records []Record
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Fatal(err)
	}

	dataMutex.Lock()
	for _, record := range records {
		recordData[record.ID] = record
	}
	dataMutex.Unlock()
	return nil
}

func formatDate(gameData map[int]Game) map[int]Game {
	formattedData := make(map[int]Game, len(gameData))
	for i, game := range gameData {
		gameCopy := game
		gameTime := gameCopy.StartDate

		t, err := time.Parse(time.RFC3339, gameTime)
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			panic(err)
		}

		inPT := t.In(loc)
		formattedDate := inPT.Format("01/02 03:04 PM")

		gameCopy.StartDate = formattedDate
		formattedData[i] = gameCopy
	}
	return formattedData
}

func ByConference(gameData map[int]Game) map[string][]Game {
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

func ByFeatured(gameData map[int]Game) []Game {
	featTeams := []Game{}
	teams := []string{"Arizona Wildcats", "Arizona State Sun Devils", "California Golden Bears", "Oregon Ducks", "Oregon State Beavers",
		"Standford Cardinal"}
	for i := range gameData {
		for j := range teams {
			if gameData[i].AwayTeam.Name == teams[j] || gameData[i].HomeTeam.Name == teams[j] {
				featTeams = append(featTeams, gameData[i])
			}
		}
	}
	return featTeams
}

func ByConferenceStandings(recordData map[int]Record) map[string][]Record {
	sortedConf := make(map[string][]Record)
	conferences := []string{"ACC", "American Athletic", "Big 12", "Big Ten",
		"Conference USA", "FBS Independent", "Mid-American",
		"Mountain West", "Pac-12", "SEC", "Sun Belt"}
	for i := range recordData {
		for j := range conferences {
			if recordData[i].Conference == conferences[j] {
				sortedConf[conferences[j]] = append(sortedConf[conferences[j]], recordData[i])
			}
		}
	}
	return sortedConf
}

func handleGames(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := gamesTpl.Execute(w, ByConference(gameData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleFeatured(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := featuredTpl.Execute(w, ByFeatured(gameData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleStandings(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	w.Header().Set("Content-Type", "text/html")
	err := standingsTpl.Execute(w, ByConferenceStandings(recordData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

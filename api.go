package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

func readSample(filename string) []byte {
	data, err := os.ReadFile("data/" + filename)
  if err != nil {
    fmt.Println(err)
  }
	return data
}

func unmarshalSampleGames(data []byte) error{
  err := json.Unmarshal(data, &games)
  if err != nil {
    fmt.Println(err)
		return err
  }
	return nil
}

func loadSampleGames(games []Game) {
  dataMutex.Lock()
  for _, game := range games {
    gameData[game.ID] = game
  }
  dataMutex.Unlock()
}

func unmarshalSampleRecords(data []byte) error{
  err := json.Unmarshal(data, &records)
  if err != nil {
    fmt.Println(err)
		return err
  }
	return nil
}

func loadSampleRecords(records []Record) {
  dataMutex.Lock()
  for _, record := range records {
    recordData[record.ID] = record
  }
  dataMutex.Unlock()
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

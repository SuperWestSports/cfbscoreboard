package main

import (
	"math"
	"sort"
	"time"
)

func formatDate() map[int]Game {
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


func ByConference() map[string][]Game {
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

	status := map[string]int{"in_progress": 0, "scheduled": 1, "completed": 2}
	for _, conference := range conferences {
		game := sortedConf[conference]
		sort.Slice(game, func(i, j int) bool {
			return status[game[i].Status] < status[game[j].Status]
		})
		sortedConf[conference] = game
	}

	return sortedConf
}

func ByFeatured() []Game {
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

	status := map[string]int{"in_progress": 0, "scheduled": 1, "completed": 2}
	sort.Slice(featTeams, func(i, j int) bool {
		return status[featTeams[i].Status] < status[featTeams[j].Status]
	})
	
	return featTeams
}

func ByConferenceStandings() map[string][]Record {
	sortedConf := make(map[string][]Record)
	conferences := []string{"ACC", "American Athletic", "Big 12", "Big Ten",
		"Conference USA", "FBS Independent", "Mid-American",
		"Mountain West", "Pac-12", "SEC", "Sun Belt"}

	for _, record := range recordData {
		for _, conference := range conferences {
			if record.Conference == conference {
				sortedConf[conference] = append(sortedConf[conference], record)
			}
		}
	}

	for _, conference := range conferences {
		records := sortedConf[conference]
		sort.Slice(records, func(i, j int) bool {
			return records[i].Total.Pct > records[j].Total.Pct
		})
		sortedConf[conference] = records
	}

	return sortedConf
}

func addPct(records []Record) []Record {
	data := make([]Record, len(records))
	for i, record := range records {
		recordCopy := record
		recordPercentage := (float64(recordCopy.Total.Wins) / float64(recordCopy.Total.Games)) 
		multiplier := math.Pow(10, 3) 
		roundedPercentage := math.Round(recordPercentage * multiplier) / multiplier
		recordCopy.Total.Pct = roundedPercentage
		data[i] = recordCopy
	}
	return data
}

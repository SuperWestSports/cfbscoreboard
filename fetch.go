package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	defer dataMutex.Unlock()
	for _, game := range games {
		gameData[game.ID] = game
	}

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

	records = addPct(records)

	dataMutex.Lock()
	defer dataMutex.Unlock()
	for _, record := range records {
		recordData[record.ID] = record
	}

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
	defer dataMutex.Unlock()
	for _, game := range games {
		gameData[game.ID] = game
	}

	return nil
}

func fetchSampleRecordsData(filename string) error {
	data, err := os.ReadFile("data/" + filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var records []Record
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Fatal(err)
	}

	records = addPct(records)

	dataMutex.Lock()
	defer dataMutex.Unlock()

	for _, record := range records {
		recordData[record.ID] = record
	}

	return nil
}


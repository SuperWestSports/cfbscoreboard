package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)


func collegeFootballData() {
	req := formatRequest()
	responseBytes := sendRequest(req)
	writeToJson(responseBytes)
}

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

func sendRequest(req *http.Request) []byte {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}
	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading HTTP resonse body: %v", err)
	}
	return responseBytes
}

func writeToJson(responseBytes []byte) {
	err := os.WriteFile("data/2023week1.json", responseBytes, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

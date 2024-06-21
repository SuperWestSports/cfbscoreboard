package main

import (
	"fmt"
)

func apiSettings() string {
	section := "games"
	year := "2023"
	week := "1"
	// KEY := "RMy62JITIczdOcIcgpVpLhfOsl4BlOFvLWsW/NGM/ZgiCcbL3bRK7JnbISToCImy"
	return fmt.Sprintf("https://api.collegefootballdata.com/%s?year=%s&week=%s", section, year, week)
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


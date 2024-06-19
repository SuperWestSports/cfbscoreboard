package main

import (
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
	//collegeFootballData()
	createTemplate()
}

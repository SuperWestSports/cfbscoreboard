package main

// Scoreboard

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


// Standings
type Record struct {
	Year	          int `json:"year"`
	ID	            int `json:"teamId"`
	Team	          string `json:"team"`
	Conference	    string `json:"conference"`
	Division	      string `json:"division"`
	Total           Games `json:"total"`
	ConferenceGames Games `json:"conferenceGames"`
	HomeGames	      Games `json:"homeGames"`
	AwayGames       Games `json:"awayGames"`
	
}

type Games struct {
	Games	  int `json:"games"`
	Wins	  int `json:"wins"`
	Losses	int `json:"losses"`
	Ties	  int `json:"ties"`
}


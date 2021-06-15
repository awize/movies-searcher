package model

type Movie struct {
	ID       int    `json:"imdbID"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Budget   int    `json:"budget"`
	Domgross int    `json:"domgross"`
}

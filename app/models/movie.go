package models

type Movie struct {
	Year   int     `json:"year"`
	Title  string  `json:"title"`
	Plot   string  `json:"plot"`
	Rating float64 `json:"rating"`
}



package model

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	EnglishName string `json:"english_name"`
	Name        string `json:"name"`
}

type Movie struct {
	ID               int        `json:"id"`
	Budget           int        `json:"budget"`
	Genres           []Genre    `json:"genres"`
	OriginalLanguage string     `json:"original_language"`
	OriginalTitle    string     `json:"original_title"`
	Overview         string     `json:"overview"`
	PosterPath       string     `json:"poster_path"`
	ReleaseDate      string     `json:"release_date"`
	Revenue          int        `json:"revenue"`
	SpokenLanguages  []Language `json:"spoken_languages"`
	Status           string     `json:"status"`
	Title            string     `json:"title"`
	VoteAverage      float64    `json:"vote_average"`
	VoteCount        int        `json:"vote_count"`
}

// type MovieDBResponse struct {
// 	ID               int                 `json:"ID"`
// 	Budget           int                 `json:"budget"`
// 	Genres           []map[string]string `json:"genres"`
// 	OriginalLanguage string              `json:"original_language"`
// 	OriginalTitle    string              `json:"original_title"`
// 	Overview         string              `json:"overview"`
// 	PosterPath       string              `json:"poster_path"`
// 	ReleaseDate      string              `json:"release_date"`
// 	Revenue          int                 `json:"revenue"`
// 	SpokenLanguages  []map[string]string `json:"spoken_languages"`
// 	Status           string              `json:"status"`
// 	Title            string              `json:"title"`
// 	VoteAverage      float64             `json:"vote_average"`
// 	VoteCount        int                 `json:"vote_count"`
// }

// genres 28:action-

// {
// 	"budget": 200000000,
//     "genres": [
//         {
//             "id": 28,
//             "name": "Action"
//         },
//         {
//             "id": 12,
//             "name": "Adventure"
//         },
//         {
//             "id": 35,
//             "name": "Comedy"
//         }
//     ],
// 	"id": 384018,
// 	"original_language": "en",
// 	"original_title": "Fast & Furious Presents: Hobbs & Shaw",
// 	"overview": "Ever since US Diplomatic Security Service Agent Hobbs and lawless outcast Shaw first faced off, they just have traded smack talk and body blows. But when cyber-genetically enhanced anarchist Brixton's ruthless actions threaten the future of humanity, they join forces to defeat him.",
// 	"poster_path": "/qRyy2UmjC5ur9bDi3kpNNRCc5nc.jpg",
// 	"release_date": "2019-08-01",
// 	"revenue": 760098996,
// 	"runtime": 137,
// 	"spoken_languages": [
// 			"English",

// 			"Russian",

// 			"Samoan",

// 	],
// 	"status": "Released",
// 	"title": "Fast & Furious Presents: Hobbs & Shaw",
// 	"vote_average": 6.9,
// 	"vote_count": 4973
// }

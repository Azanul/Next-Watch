package model

type Movie struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Picture  *string      `json:"picture"`
	Director Attribute    `json:"director"`
	Actors   []*Attribute `json:"actors"`
	Genres   []*Attribute `json:"genres"`
	Rating   int          `json:"rating"`
}

type Attribute struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

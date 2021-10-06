package main

import (
	"encoding/json"
	"net/http"
)

type Munro struct {
	Name       string `json:"name"`
	Area       string `json:"area"`
	Time       string `json:"time"`
	Difficulty string `json:"difficulty"`
	Height     string `json:"height"`
}

type munroHandlers struct {
	store map[string]Munro
}

func (h *munroHandlers) get(w http.ResponseWriter, r *http.Request) {
	munros := make([]Munro, len(h.store))

	i := 0
	for _, munro := range h.store {
		munros[i] = munro
		i++
	}

	jsonBytes, err := json.Marshal(munros)
	if err != nil {
		//TODO
	}

	w.Write(jsonBytes)
}

func newMunroHandlers() *munroHandlers {
	return &munroHandlers{
		store: map[string]Munro{
			"id1": Munro{
				Name:       "Schiehallion",
				Area:       "Aberfoyle",
				Time:       "6 hrs",
				Difficulty: "7/10",
				Height:     "915m",
			},
		},
	}
}

func main() {
	munroHandlers := newMunroHandlers()
	http.HandleFunc("/munros", munroHandlers.get)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//complete up to 6 mins of vid

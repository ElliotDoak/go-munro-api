package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Munro struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Area       string `json:"area"`
	Time       string `json:"time"`
	Difficulty string `json:"difficulty"`
	Height     string `json:"height"`
}

type munroHandlers struct {
	sync.Mutex
	store map[string]Munro
}

func (h *munroHandlers) get(w http.ResponseWriter, r *http.Request) {
	munros := make([]Munro, len(h.store))

	h.Lock()
	i := 0
	for _, munro := range h.store {
		munros[i] = munro
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(munros)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *munroHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var munro Munro
	json.Unmarshal(bodyBytes, &munro)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	munro.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[munro.Name] = munro
	defer h.Unlock()
}

func (h *munroHandlers) munros(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
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
	http.HandleFunc("/munros", munroHandlers.munros)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//completed up to 24mins of vid

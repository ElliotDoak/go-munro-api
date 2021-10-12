package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

func (h *munroHandlers) getRandomMunro(w http.ResponseWriter, R *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Add("location", fmt.Sprintf("/munros/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *munroHandlers) getMunro(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomMunro(w, r)
		return
	}

	h.Lock()
	munro, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(munro)
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
		return
	}

	//ct := r.Header.Get("content-type")
	// if ct != "application/json" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
	// }
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
			// "id1": Munro{
			// 	ID:         "123",
			// 	Name:       "Schiehallion",
			// 	Area:       "Aberfoyle",
			// 	Time:       "6 hrs",
			// 	Difficulty: "7/10",
			// 	Height:     "915m",
			// },
		},
	}
}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}
	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorised"))
		return
	}
	w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	munroHandlers := newMunroHandlers()
	http.HandleFunc("/munros", munroHandlers.munros)
	http.HandleFunc("/munros/", munroHandlers.getMunro)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

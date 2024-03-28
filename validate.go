package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirps struct {
	Body string `json:"body"`
}
type resCleanMsg struct {
	CleanedBody string `json:"cleaned_body"`
}
type resErr struct {
	Error string `json:"error"`
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// decode request
	decoder := json.NewDecoder(r.Body)
	chirp := chirps{}
	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	// check length of the chirp
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	w.WriteHeader(200)
	// validate profane
	newmsg := replaceBadWords(chirp.Body)
	res := resCleanMsg{CleanedBody: newmsg}
	dat, _ := json.Marshal(res)
	w.Write([]byte(dat))
}

func replaceBadWords(msg string) (newmsg string) {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(msg)
	for i, word := range words {
		for _, profane := range profane_words {
			if strings.ToLower(word) == profane {
				words[i] = "****"
			}
		}
	}
	newmsg = strings.Join(words, " ")
	return
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	w.WriteHeader(code)
	res := resErr{Error: msg}
	dat, _ := json.Marshal(res)
	w.Write([]byte(dat))
}

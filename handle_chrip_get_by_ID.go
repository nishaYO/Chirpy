package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerChirpRetrieveByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := getChirpID(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func getChirpID(r *http.Request) (int, error) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	lastParam := parts[len(parts)-1]
	chirpID, err := strconv.Atoi(lastParam)
	if err != nil {
		return 0, fmt.Errorf("Chirp ID %v is not a valid integer", lastParam)
	}
	return chirpID, nil
}

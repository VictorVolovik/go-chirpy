package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	replacer := "****"
	cleaned := replaceWords(params.Body, badWords, replacer)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})

}

func replaceWords(str string, wordsToReplace map[string]struct{}, replacer string) string {
	strWords := strings.Fields(str)

	for i, word := range strWords {
		wordLowercased := strings.ToLower(word)
		if _, ok := wordsToReplace[wordLowercased]; ok {
			strWords[i] = replacer
		}
	}

	joined := strings.Join(strWords, " ")
	return joined
}

package main

import (
	"encoding/json"
	"net/http"
	"time"
	"errors"
	"strings"

	"github.com/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct{
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	UserID 		uuid.UUID 	`json:"user_id"`	
	Body 		string 		`json:"body"`
}


func (cfg *apiConfig) handlerCreateChirp (w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body 	string		`json:"body"`
		UserID	uuid.UUID	`json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.dbPtr.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: 		chirp.ID,
		CreatedAt: 	chirp.CreatedAt,
		UpdatedAt: 	chirp.UpdatedAt,
		Body: 		chirp.Body,
		UserID: 	chirp.UserID,
	})

}

func validateChirp(body string) (string, error){
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words{
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok{
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}


func (cfg *apiConfig) handlerGetAllChrips(w http.ResponseWriter, r *http.Request){
	chirpsArray := []Chirp{}
	chirpsInDB, err := cfg.dbPtr.GetAllChrips(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error with GetAllChrips Query", err)
		return
	}

	for _, c := range chirpsInDB {
		chrp := Chirp{
			ID: 		c.ID,
			CreatedAt: 	c.CreatedAt,
			UpdatedAt: 	c.UpdatedAt,
			UserID: 	c.UserID,
			Body: 	c.Body,
		}

		chirpsArray = append(chirpsArray, chrp)
	}

	respondWithJSON(w, http.StatusOK, chirpsArray)

}

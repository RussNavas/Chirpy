package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"sort"

	"github.com/Chirpy/internal/auth"
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
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldnt get bearer token", err)
		return
	}

	usrId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt validate jwt", fmt.Errorf("The err: %v", err))
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.dbPtr.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: usrId,
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
	sortby := r.URL.Query().Get("sort")
	authorIDStr := r.URL.Query().Get("author_id")
	if authorIDStr != ""{
		parsedID , err := uuid.Parse(authorIDStr)
		if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error parsing author id", err)
		return			
		}
		userChrips, err := cfg.dbPtr.GetUserChirps(r.Context(), parsedID)
		if err != nil{
			respondWithError(w, http.StatusInternalServerError, "Error getting authors chirps", err)
		return				
		}

		cArray := []Chirp{}

		for _, c := range userChrips{
		chrp := Chirp{
			ID: 		c.ID,
			CreatedAt: 	c.CreatedAt,
			UpdatedAt: 	c.UpdatedAt,
			UserID: 	c.UserID,
			Body: 	c.Body,
		}

		cArray = append(cArray, chrp)
	}

	if sortby == "desc"{
		sort.Slice(cArray, func(i, j int) bool {
			return cArray[i].CreatedAt.After(cArray[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, cArray)
	return
	}


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

	if sortby == "desc"{
		sort.Slice(chirpsArray, func(i, j int) bool {
			return chirpsArray[i].CreatedAt.After(chirpsArray[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsArray)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request){
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == ""{
		respondWithError(w, 404, "Chirp ID Not Found", errors.New(""))
		return
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, 404, "Invalid chirp ID str", err)
		return
	}

	c, err := cfg.dbPtr.GetChirpByChirpID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Couldn't get chirp from DB", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body: c.Body,
		UserID: c.UserID,
	} )
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/Chirpy/internal/auth"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type Parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode params", err)
		return
	}

	
	usr, err := cfg.dbPtr.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt get user by email", err)
		return
	}

	status, err := auth.CheckPasswordHash(params.Password, usr.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error checking password hash", err)
		return
	}

	if status != true{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email: usr.Email,
	})
}


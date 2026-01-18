package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/Chirpy/internal/auth"
	"github.com/google/uuid"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type Parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int	`json:"expires_in_seconds"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode params", err)
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600{
		params.ExpiresInSeconds = 3600
	}

	usr, err := cfg.dbPtr.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt get user by email", err)
		return
	}

	jwt, err := auth.MakeJWT(usr.ID, cfg.secret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "error making JWT", err)
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

	type LoginResponse struct{
		ID 			uuid.UUID 	`json:"id"`
		CreatedAt 	time.Time 	`json:"created_at"`
		UpdatedAT 	time.Time 	`json:"updated_at"`
		Email		string		`json:"email"`
		Token		string 		`json:"token"`
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAT: usr.UpdatedAt,
		Email: usr.Email,
		Token: jwt,
	})
}


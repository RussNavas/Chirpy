package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Chirpy/internal/auth"
	"github.com/Chirpy/internal/database"
	"github.com/google/uuid"
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

	jwt, err := auth.MakeJWT(usr.ID, cfg.secret, time.Hour)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "error making JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "error making Refresh Token", err)
	}

	cfg.dbPtr.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: usr.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24*60),
	})

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
		ID 				uuid.UUID 	`json:"id"`
		CreatedAt 		time.Time 	`json:"created_at"`
		UpdatedAT 		time.Time 	`json:"updated_at"`
		Email			string		`json:"email"`
		IsChirpyRed		bool		`json:"is_chirpy_red"`
		Token			string 		`json:"token"`
		RefreshToken 	string		`json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAT: usr.UpdatedAt,
		Email: usr.Email,
		IsChirpyRed: usr.IsChirpyRed,
		Token: jwt,
		RefreshToken: refreshToken,
	})
}


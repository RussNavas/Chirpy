package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){
	authHeader := r.Header.Get("Authorization")
	if authHeader == ""{
		respondWithError(w, 401, "No auth header found", fmt.Errorf("no auth header found"))
		return
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) < 2 || splitAuth[0] != "Bearer"{
		respondWithError(w, 401, "malformed authorization header", fmt.Errorf("malformed authorization header"))
	}
	token := splitAuth[1]


	dbToken, err := cfg.dbPtr.GetRefreshTokenByToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Error getting token form db", err)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()){
		respondWithError(w, 401, "Token Expired", fmt.Errorf("Token Expired"))
		return
	}
	if dbToken.RevokedAt.Valid{
		respondWithError(w, 401, "Token Expired", fmt.Errorf("Token Expired"))
		return
	}
	type payload struct{
		Token string `json:"token"`
	}

	newAccessToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret, time.Hour)
	if err != nil{
		respondWithError(w, 401, "Couldnt make new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, payload{
		Token: newAccessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){
	authHeader := r.Header.Get("Authorization")
	if authHeader == ""{
		respondWithError(w, 401, "No auth header found", fmt.Errorf("no auth header found"))
		return
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) < 2 || splitAuth[0] != "Bearer"{
		respondWithError(w, 401, "malformed authorization header", fmt.Errorf("malformed authorization header"))
		return
	}
	token := splitAuth[1]
	dbToken, err := cfg.dbPtr.GetRefreshTokenByToken(r.Context(), token)
	if err != nil{
		respondWithError(w, 401, "Error getting token from db", err)
		return
	}

	_, err = cfg.dbPtr.UpdateRevokeToken(r.Context(), dbToken.Token)
	if err != nil{
		respondWithError(w, 401, "Error updating token from db", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

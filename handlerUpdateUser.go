package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chirpy/internal/auth"
	"github.com/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request){
	type Parameters struct{
		Password		string		`json:"password"`
		Email			string		`json:"email"`
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "could update user error decoding into params", err)
		return
	}


	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "error hashing password", err)
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

	updatedUsr, err := cfg.dbPtr.UpdateUserPasswordandEmailByID(r.Context(), database.UpdateUserPasswordandEmailByIDParams{
		Email: params.Email,
		HashedPassword: hashedPwd,
		ID: usrId,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt validate jwt", fmt.Errorf("The err: %v", err))
		return
	}

	type Response struct{
		ID 				uuid.UUID 	`json:"id"`
		Created_at		time.Time 	`json:"created_at"`
		Updated_at		time.Time 	`json:"updated_at"`
		Email			string		`json:"email"`
	}

	respondWithJSON(w, http.StatusOK, Response{
		ID: updatedUsr.ID,
		Created_at: updatedUsr.CreatedAt,
		Updated_at: updatedUsr.UpdatedAt,
		Email:		updatedUsr.Email,
	})


}

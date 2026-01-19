package main

import (
	"fmt"
	"net/http"


	"github.com/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request){

	// parse for chirp ID in URL
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == ""{
		fmt.Printf("problem with r.pathValue chirp ID\n")
		respondWithError(w, 404, "Chirp ID Not Found", fmt.Errorf(""))
		return
	}

	// convert to uuid
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("problem with parsing chirpIDStr\n")
		respondWithError(w, 404, "Invalid chirp ID str", err)
		return
	}

	// search for token
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		fmt.Printf("problem with getting bearer token from header to prepare for deletion\n")
		respondWithError(w, http.StatusUnauthorized, "Couldnt get bearer token", err)
		return
	}

	// validate the token
	usrID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		fmt.Printf("problem with validating JWT to prepare for deletion\n")
		respondWithError(w, http.StatusUnauthorized, "couldnt validate jwt", fmt.Errorf("The err: %v", err))
		return
	}

	// get the chirp 
	chrp, err := cfg.dbPtr.GetChirpByChirpID(r.Context(), chirpID)
	if err != nil {
		fmt.Printf("problem with getting chirp by chirp ID\n")
		respondWithError(w, http.StatusNotFound, "couldnt get chirp by ID", fmt.Errorf("The err: %v", err))
		return
	}

	// make sure it is the correct chirp
	if chrp.UserID != usrID{
		fmt.Printf("userID mismatch\n")
		respondWithError(w, http.StatusForbidden, "UserID does not belong to chirp", fmt.Errorf("The err: %v", err))
		return
	}

	// delete it
	err = cfg.dbPtr.DeleteChirp(r.Context(), chirpID)
	if err != nil{
		fmt.Printf("problem with deleting the chirp\n")
		respondWithError(w, http.StatusInternalServerError, "Couldnt delete chirp", err)
		return
	}

	// confirm the deletion
	fmt.Printf("Success!\n")
	w.WriteHeader(204)
}

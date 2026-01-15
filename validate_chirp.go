package main

import(
	"net/http"
	"encoding/json"
)

func handlerValidate(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}

	type returnVals struct{
		Valid bool `json:"valid"`
	}

	type containedProfanity struct{
		Cleaned_Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	strs, neededCleaning := filterProfanity(params.Body)
	if neededCleaning == true{
		respondWithJSON(w, http.StatusOK, containedProfanity{
			Cleaned_Body: strs,
		})
		return
	}
	
	respondWithJSON(w, http.StatusOK, containedProfanity{
		Cleaned_Body: params.Body,
	})

}

package main


import(
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks (w http.ResponseWriter, r *http.Request){
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		w.WriteHeader(http.StatusUnauthorized)
	}

	if apiKey != cfg.polkaKey{
		w.WriteHeader(http.StatusUnauthorized)
	}

	type Data struct{
		UserID	string	`json:"user_id"`
	}
	
	type Parameters struct{
		Event	string	`json:"event"`
		Data	Data	`json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}

	err = decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "problem decoding", err)
		return
	}

	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userIDstr := params.Data.UserID
	convertedUserID, err := uuid.Parse(userIDstr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem parsing to uuid", err)
		return
	}

	_, err = cfg.dbPtr.UpgradeToRed(r.Context(), convertedUserID)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Couldnt find user to upgrade", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

package main

import(
	"github.com/google/uuid"
	"encoding/json"
	"time"
	"net/http"
)


type User struct{
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	Email		string		`json:"email"`
}

func handlerCreateUser(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	usr := User{}
	err := decoder.Decode(&usr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user", err)
		return
	}

	respondWithJSON(w, 201, usr)
}

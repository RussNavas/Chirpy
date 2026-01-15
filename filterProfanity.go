package main

import(
	"strings"
)
// postgres://russellnavas:@localhost:5432/chirpy
func filterProfanity(str string) (string, bool){
	neededCleaning := false
	// inspect message for profanity
	strs := strings.Split(str, " ")
	for i, s := range strs{
		s_lower := strings.ToLower(s)
		switch s_lower{
		case "kerfuffle":
			neededCleaning = true
			strs[i] = "****"
		case "sharbert":
			neededCleaning = true
			strs[i] = "****"
		case "fornax":
			neededCleaning = true
			strs[i] = "****"
		}
	}
	// package cleaned message
	cleaned := strings.Join(strs, " ")
	return cleaned, neededCleaning
}


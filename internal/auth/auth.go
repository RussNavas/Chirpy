package auth

import(
	"github.com/alexedwards/argon2id"
	"fmt"
)

func HashPassword(password string) (string, error){
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil{
		return "", fmt.Errorf("couldnt hash password: %v", err)
	}

	return hash, nil	
}

func CheckPasswordHash(password, hash string) (bool, error){
	status, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil{
		return false, fmt.Errorf("error comparing pass and hash: %v", err)
	}

	return status, nil
}

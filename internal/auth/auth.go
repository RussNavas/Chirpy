package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"crypto/rand" 
	"encoding/hex"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration)(string, error){
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject: userID.String(),
	})

	str, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}
	return str, nil
}

func ValidateJWT(tokenString, tokenSecret string)(uuid.UUID, error){
	claim := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token)(interface{}, error){
		return []byte(tokenSecret), nil
	}  )
	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid{
		return uuid.UUID{}, fmt.Errorf("Token is invalid")
	}

	id := claim.Subject
	uuID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuID, nil
}


func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == ""{
		return "", fmt.Errorf("no authorization header included")
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) < 2 || splitAuth[0] != "Bearer"{
		return "", fmt.Errorf("malformed authorization header")
	}
	return splitAuth[1], nil
}

func MakeRefreshToken() (string, error){
	arr := make([]byte, 32)
	_, err := rand.Read(arr)
	if err != nil{
		return "", fmt.Errorf("couldnt rand.Read into arr: %v", err)
	}

	numStr := hex.EncodeToString(arr)

	return numStr, nil
}

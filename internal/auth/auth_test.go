package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashAndCompare(t *testing.T){
	password := "RubberDucky123"
	hashed, err := HashPassword(password)
	if err != nil{
		t.Errorf("Problem hashing password for test: %v", err)
		return
	}
	status, err := CheckPasswordHash(password, hashed)
	if err != nil{
		t.Errorf("Problem comparing password and hash: %v", err)
	}

	if status == true{
		t.Logf("SUCCESS! password was hashed and checked")
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "super_safe_secret_key"
	userID := uuid.New()

	// Case 1: Valid Token
	// Create a token that expires in 1 hour
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT unexpected error: %v", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT unexpected error: %v", err)
	}
	if validatedID != userID {
		t.Errorf("expected user ID %v, got %v", userID, validatedID)
	}

	// Case 2: Expired Token
	// Create a token that expired 1 second ago
	expiredToken, err := MakeJWT(userID, secret, -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT unexpected error: %v", err)
	}

	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("ValidateJWT expected error for expired token, got nil")
	}

	// Case 3: Wrong Secret
	// Use the valid token from Case 1, but try to validate it with a different key
	wrongSecret := "malicious_key"
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT expected error for wrong secret, got nil")
	}
}

package auth

import (
	"testing"
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

package service

import (
	"github.com/golang-jwt/jwt"
	"math/rand"
	"testing"
)

const signingKeyTest = "j370sdfs34472fshvlruso043275fhka"

func TestCheckJWT(t *testing.T) {

	testUserID := rand.Int()
	testClaims := &tokenClaims{
		UserId: testUserID,
	}
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	testTokenString, err := testToken.SignedString([]byte(signingKeyTest))
	if err != nil {
		t.Fatalf("Error creating test token: %v", err)
	}
	t.Run("ValidToken", func(t *testing.T) {
		resultUserID, err := CheckJWT(testTokenString)
		if err != nil {
			t.Fatalf("CheckJWT returned an error: %v", err)
		}

		if resultUserID != testUserID {
			t.Errorf("Expected user ID %v, got %v", testUserID, resultUserID)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		invalidTokenString := "invalid_token_string"
		_, err = CheckJWT(invalidTokenString)
		if err == nil {
			t.Error("CheckJWT did not return an error for an invalid token")
		}
	})
}

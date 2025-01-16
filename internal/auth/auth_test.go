package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing_ValidatesCorrectly(t *testing.T) {
	cases := []struct {
		password string
	}{
		{
			password: "qwerty123",
		},
		{
			password: "superman",
		},
	}

	for _, c := range cases {
		actual, err := HashPassword(c.password)
		if err != nil {
			t.Errorf("should hash password, %s", err)
			return
		}
		expected := CheckPasswordHash(c.password, actual)
		if expected != nil {
			t.Errorf("should match hashes")
			return
		}
	}
}

func TestPasswordHashing_RejectsIncorrect(t *testing.T) {
	correctPassword := "Sup3rS3cr3tPa$$w0rd"
	wrongPassword := "tooeasy"

	actual, err := HashPassword(correctPassword)
	if err != nil {
		t.Errorf("should hash password, %s", err)
		return
	}
	expected := CheckPasswordHash(wrongPassword, actual)
	if expected == nil {
		t.Errorf("should reject incorrect password")
		return
	}

}

func TestValidateJWT_Valid(t *testing.T) {
	secret := "testsecret"
	expiration := time.Hour
	expectedUserID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("should create user id, %s", err)
		return
	}

	jwt, err := MakeJWT(expectedUserID, secret, expiration)
	if err != nil {
		t.Errorf("should make new jwt, %s", err)
		return
	}

	actualUserID, err := ValidateJWT(jwt, secret)
	if err != nil {
		t.Errorf("should validate jwt, %s", err)
		return
	}

	if expectedUserID != actualUserID {
		t.Errorf("should match users")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	secret := "testsecret"
	expiration := time.Hour
	expectedUserID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("should create user id, %s", err)
		return
	}

	jwt, err := MakeJWT(expectedUserID, secret, expiration)
	if err != nil {
		t.Errorf("should make new jwt, %s", err)
		return
	}

	_, err = ValidateJWT(jwt, "someothersecret")
	if err == nil {
		t.Errorf("should error on secret mismatch")
		return
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	secret := "testsecret"
	expiration := -time.Hour
	expectedUserID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("should create user id, %s", err)
		return
	}

	jwt, err := MakeJWT(expectedUserID, secret, expiration)
	if err != nil {
		t.Errorf("should make new jwt, %s", err)
		return
	}

	_, err = ValidateJWT(jwt, secret)
	if err == nil {
		t.Errorf("should error as expired")
		return
	}
}

func TestGetBearerToken_Exists(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer TOKEN_STRING")
	expected := "TOKEN_STRING"

	actual, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("should get bearer token")
	}

	if actual != expected {
		t.Errorf("should match")
	}
}

func TestGetBearerToken_NoAuthHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Errorf("should error as auth headers not found")
	}
}

func TestGetBearerToken_MalformeddHeader(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Errorf("should error as malformed auth header")
	}
}

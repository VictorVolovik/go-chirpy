package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		input string
	}{
		{
			input: "qwerty123",
		},
		{
			input: "superman",
		},
	}

	for _, c := range cases {
		actual, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("unable to hash password, %s", err)
		}
		expected := CheckPasswordHash(c.input, actual)
		if expected != nil {
			t.Errorf("hash doesn't match")
		}
	}
}

package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	header1 := http.Header{}
	header1.Add("Authorization", "Bearer Joseph")

	cases := []struct {
		input          http.Header
		expectedString string
		expectedErr    error
	}{
		{
			input:          header1,
			expectedString: "Joseph",
			expectedErr:    nil,
		},
		{
			input:          http.Header{},
			expectedString: "",
			expectedErr:    fmt.Errorf("authorization header does not exist"),
		},
	}

	for _, c := range cases {
		actual, err := GetBearerToken(c.input)
		if err == nil {
			if c.expectedErr != nil {
				t.Errorf("Expected %v error and was given none", c.expectedErr)
				t.Fail()
			}
		} else {
			if c.expectedErr == nil {
				t.Errorf("Expected no error and was given %v", err)
				t.Fail()
			}
		}
		if actual != c.expectedString {
			t.Errorf("Expected %v, was given %v", c.expectedString, actual)
			t.Fail()
		}
	}
}

package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTToken(t *testing.T) {
	uuid, _ := uuid.Parse("(00|_ &|_|`/")

	jwt, err := MakeJWT(uuid, "53[|`[-+ |-|^<l<3|2[v]@/V", 3*time.Minute)
	if err != nil {
		t.Logf("Unable to create JWT: %v", err)
		t.Fail()
	}

	jwtUuid, err := ValidateJWT(jwt, "53[|`[-+ |-|^<l<3|2[v]@/V")
	if err != nil {
		t.Logf("Unable to validate JWT: %v", err)
		t.Fail()
	}

	if jwtUuid != uuid {
		t.Logf("uuids are not the same, expected: %v got: %v", uuid, jwtUuid)
		t.Fail()
	}
}

package authgoblue_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"authgoblue/claims"
)

func TestClaimsJSONMarshal(t *testing.T) {

	input := claims.Claims{

		UserID: "user-123",

		Username: "alice",

		Email: "alice@example.com",

		Role: "admin",

		Permissions: []string{
			"read",
			"write",
		},

		TokenType: claims.TokenTypeAccess,

		Issuer: "auth-service",

		IssuedAt: 1000,

		ExpiresAt: 2000,
	}

	data, err :=
		json.Marshal(input)

	if err != nil {
		t.Fatal(err)
	}

	var output claims.Claims

	err =
		json.Unmarshal(
			data,
			&output,
		)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(
		input,
		output,
	) {

		t.Fatalf(
			"claims mismatch\ninput=%v\noutput=%v",
			input,
			output,
		)
	}

}

func TestClaimsOptionalFields(t *testing.T) {

	input := claims.Claims{

		UserID: "user-123",
	}

	data, err :=
		json.Marshal(input)

	if err != nil {
		t.Fatal(err)
	}

	result := string(data)

	if result == "" {

		t.Fatal(
			"expected json output",
		)
	}

	if result != `{"user_id":"user-123"}` {

		t.Fatalf(
			"unexpected json: %s",
			result,
		)
	}

}

func TestClaimsEmptyPermissions(t *testing.T) {

	c := claims.Claims{

		UserID: "user-123",

		Permissions: []string{},
	}

	if len(c.Permissions) != 0 {

		t.Fatal(
			"expected empty permissions",
		)
	}

}

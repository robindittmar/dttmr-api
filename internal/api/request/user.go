package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserPayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func DecodeCreateUser(r *http.Request) (CreateUserPayload, error) {
	var payload CreateUserPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return payload, fmt.Errorf("error decoding create user payload: %w", err)
	}

	return payload, nil
}

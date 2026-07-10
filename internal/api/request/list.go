package request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateListPayload struct {
	Name    string   `json:"name"`
	UserIDs []string `json:"user_ids"`
}

func DecodeCreateList(r *http.Request) (CreateListPayload, error) {
	var payload CreateListPayload

	r.Body = http.MaxBytesReader(nil, r.Body, 1024*64)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return payload, fmt.Errorf("error decoding create list payload: %w", err)
	}

	return payload, nil
}

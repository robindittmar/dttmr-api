package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ApiResponse struct {
	Version int    `json:"version"`
	Path    string `json:"path"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := ApiResponse{Path: "/v1", Version: 1}

		b, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

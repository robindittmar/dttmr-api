package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ApiResponse struct {
	Method     string            `json:"method"`
	Url        string            `json:"url"`
	Proto      string            `json:"proto"`
	Header     map[string]string `json:"header"`
	Host       string            `json:"host"`
	RemoteAddr string            `json:"remoteAddr"`
	Form       map[string]string `json:"form"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var resp ApiResponse

		resp.Method = r.Method
		resp.Url = r.URL.String()
		resp.Proto = r.Proto
		resp.Header = make(map[string]string)
		resp.Host = r.Host
		resp.RemoteAddr = r.RemoteAddr
		resp.Form = make(map[string]string)

		for k, v := range r.Header {
			resp.Header[k] = v[0]
		}

		if err := r.ParseForm(); err == nil {
			for k, v := range r.Form {
				resp.Form[k] = v[0]
			}
		} else {
			log.Print(err)
		}

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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{Status: "ok"}

		b, err := json.Marshal(&resp)
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

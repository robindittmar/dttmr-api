package handler

import (
	"log"
	"net/http"

	"github.com/robindittmar/dttmr-api/internal/api/response"
)

type apiResponse struct {
	Method     string            `json:"method"`
	Url        string            `json:"url"`
	Proto      string            `json:"proto"`
	Header     map[string]string `json:"header"`
	Host       string            `json:"host"`
	RemoteAddr string            `json:"remoteAddr"`
	Form       map[string]string `json:"form"`
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

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
		log.Println(err)
	}

	response.JSON(w, http.StatusOK, resp)
}

package core

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
)

type ErrorObj struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Status 		bool   `json:"status"`
}

type FormError struct {
	Message string `json:"message"`
	Name string `json:"name"`
}

func status(code int) string {
	Status := map[int]string{
		200: "OK",
		201: "Created",
		400: "Bad request.",
		404: "Not found/nothing here.",
		500: "Internal Server Error.",
		444: "User already exists.",
		401: "Unauthorized.",
		403: "Forbidden.",
		405: "Method not allowed.",
		406: "Unacceptable data.",
		503: "Service unavailable.",
	}
	return Status[code]
}

func Error(w http.ResponseWriter, code int, description string) {
	e := ErrorObj{code, status(code), description, false}
	resp, err := json.Marshal(e)
	if err != nil {
		fmt.Fprint(w, "500: Critical Server Error")
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(resp))
	return
}

func Data(w http.ResponseWriter, data interface{}) {
	if data != nil {
		resp, err := json.Marshal(data)
		if err != nil {
			fmt.Fprint(w, "500: Critical Server Error")
			return
		}
		w.WriteHeader(http.StatusOK)
		buf := bytes.Buffer{}
		json.HTMLEscape(&buf, resp)
		w.Write(buf.Bytes())
	}else{
		w.WriteHeader(http.StatusOK)
	}
	return
}


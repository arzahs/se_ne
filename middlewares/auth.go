package middlewares

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"se_ne/models"
)

type Middleware struct {
	Session *sessions.CookieStore
}

func (mwr *Middleware) MustAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			// check in store
			user, err := models.GetUserByToken(authorizationHeader)
			if err != nil {
				BadAccessResponse(w, req, "/login/")
				return
			}
			ctx := context.WithValue(req.Context(), "User", user)
			next(w, req.WithContext(ctx))
			return
		}

		session, err := mwr.Session.Get(req, "user-session")
		if err == nil {
			value, ok := session.Values["user-token"]
			if ok && value != "" {
				user, err := models.GetUserByToken(value.(string))
				if err != nil {
					BadAccessResponse(w, req, "/login/")
					return
				}
				ctx := context.WithValue(req.Context(), "user_id", user.Id)
				next(w, req.WithContext(ctx))
				return
			}
		}

		BadAccessResponse(w, req, "/login/")
	})
}

func (mwr *Middleware) MustAnon(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			BadAccessResponse(w, req, "/")
			return
		}

		value, err := mwr.Session.Get(req, "user-session")
		if err == nil {
			value, _ := value.Values["user-token"]
			valueStr, ok := value.(string)
			if valueStr != "" || ok {
				_, err := models.GetUserByToken(valueStr)
				if err == nil {
					BadAccessResponse(w, req, "/")
					return
				}
			}
		}

		next(w, req)
	})

}

func BadAccessResponse(w http.ResponseWriter, req *http.Request, redirect string) {
	isAJAX := req.Header.Get("X-Requested-With")
	if isAJAX != "" {
		response, err := json.Marshal(&struct {
			Error  string
			Status bool
		}{Error: "Unauthorazed",
			Status: false})
		if err != nil {
			w.Write(response)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		http.Redirect(w, req, redirect, 301)
	}
}

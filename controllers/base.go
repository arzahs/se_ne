package controllers

import (
	"net/http"
	"se_ne/core"
	"github.com/gorilla/sessions"
	"se_ne/models"
)

type Controller struct{
	View *core.View
	Cfg *core.Config
	Session *sessions.CookieStore
}

type APIControllerInterface interface {
	Put(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func Login(user *models.User, store *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (string, error){
	token, _ := core.GenerateToken()
	userSession := models.Session{Token: *token, UserId: user.Id}
	err := userSession.Insert()
	if err != nil{
		return "", err
	}
	sessionStore, err := store.Get(r,"user-session")
	if err != nil{
		return "", err
	}
	sessionStore.Values["user-token"] = token
	sessionStore.Save(r, w)
	return *token, err
}




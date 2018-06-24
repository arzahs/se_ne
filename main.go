package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"se_ne/controllers"
	"se_ne/core"
	"se_ne/middlewares"
	"se_ne/models"
	"github.com/gorilla/sessions"
)

func main(){
	// init config application
	cfg, err := core.NewConfig()
	if err != nil{
		log.Printf("application run error = %s", err)
		return
	}
	err = models.InitModels(cfg.DB)
	if err != nil{
		log.Printf("db error = %s", err)
		return
	}

	// initial depend
	view := core.NewView(cfg.TemplatePath)
	sessionStore := sessions.NewCookieStore([]byte(cfg.SecretKey))
	accessMwr := middlewares.Middleware{Session: sessionStore}
	baseCtrl := controllers.Controller{Cfg: cfg, View:view, Session: sessionStore}
	pageCtrl := controllers.PageController{Controller: baseCtrl}
	userCtrl := controllers.UserController{Controller: baseCtrl}
	authCtrl := controllers.SessionController{Controller: baseCtrl}
	resetPswdCtrl := controllers.ResetRequestController{Controller: baseCtrl}
	// init router
	router := mux.NewRouter()
	// page routes
	router.HandleFunc("/", accessMwr.MustAuth(pageCtrl.Profile))
	router.HandleFunc("/login/", accessMwr.MustAnon(pageCtrl.Login))
	router.HandleFunc("/logout/", accessMwr.MustAuth(pageCtrl.Logout))
	router.HandleFunc("/registration/", accessMwr.MustAnon(pageCtrl.Registration))
	router.HandleFunc("/reset/", accessMwr.MustAnon(pageCtrl.ForgetPassword))
	router.HandleFunc("/reset/confirm/", accessMwr.MustAnon(pageCtrl.NewPassword))
	// user api
	router.HandleFunc("/api/v1/user/", accessMwr.MustAnon(userCtrl.Post)).Methods("POST")
	router.HandleFunc("/api/v1/user/", accessMwr.MustAuth(userCtrl.Get)).Methods("GET")
	router.HandleFunc("/api/v1/user/", accessMwr.MustAuth(userCtrl.Put)).Methods("PUT")
	router.HandleFunc("/api/v1/user/", accessMwr.MustAuth(userCtrl.Delete)).Methods("DELETE")
	router.HandleFunc("/api/v1/user/google/", accessMwr.MustAnon(authCtrl.GoogleStartProcess)).Methods("GET")

	// password reset api
	router.HandleFunc("/api/v1/password/", accessMwr.MustAnon(resetPswdCtrl.Post)).Methods("POST")
	router.HandleFunc("/api/v1/password/", accessMwr.MustAnon(resetPswdCtrl.Put)).Methods("PUT")
	// auth api
	router.HandleFunc("/api/v1/session/", accessMwr.MustAnon(authCtrl.Post)).Methods("POST")
	router.HandleFunc("/api/v1/session/", accessMwr.MustAuth(authCtrl.Get)).Methods("GET")
	router.HandleFunc("/api/v1/session/", accessMwr.MustAuth(authCtrl.Put)).Methods("PUT")
	router.HandleFunc("/api/v1/session/", accessMwr.MustAuth(authCtrl.Delete)).Methods("DELETE")
	router.HandleFunc("/api/v1/session/google/", accessMwr.MustAnon(authCtrl.GoogleGoogleCallback)).Methods("GET")
	// init staticfiles router
	fs := http.FileServer(http.Dir("./assets/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	// start server
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), router)
	if err != nil{
		log.Println(err.Error())
	}
}

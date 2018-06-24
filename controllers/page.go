package controllers

import (
	"net/http"
	"se_ne/core"
	"se_ne/models"
)

type PageController struct {
	Controller
}

func (pc *PageController) Login(w http.ResponseWriter, r *http.Request) {
	pc.View.Render(w, "login", nil)
}

func (pc *PageController) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := pc.Session.Get(r, "user-session")
	if err != nil {
		core.Error(w, http.StatusInternalServerError, "Session not found")
		return
	}
	value, ok := session.Values["user-token"]
	if ok && value != "" {
		models.DeleteSessionByToken(value.(string))
		delete(session.Values, "user-token")
		session.Save(r, w)
		http.Redirect(w, r, "/login/", 301)
		return
	}

	core.Error(w, http.StatusInternalServerError, "Token not found")
	return

}

func (pc *PageController) Registration(w http.ResponseWriter, r *http.Request) {
	pc.View.Render(w, "registration", nil)

}

func (pc *PageController) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	pc.View.Render(w, "reset", nil)
}

func (pc *PageController) NewPassword(w http.ResponseWriter, r *http.Request) {
	tokenList, ok := r.URL.Query()["t"]
	if !ok {
		core.Error(w, http.StatusBadRequest, "Token not found")
		return
	}
	if len(tokenList) > 0 {
		pc.View.Render(w, "new_password", struct {
			Token   string
			IsLogin bool
		}{Token: tokenList[0], IsLogin: false})
		return
	}
	core.Error(w, http.StatusBadRequest, "Token is not correct")
	return
}

func (pc *PageController) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if ok && userId != 0 {
		user, err := models.GetUserById(userId)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "User not found")
			return
		}
		pc.View.Render(w, "profile", struct {
			IsLogin      bool
			User         interface{}
			GoogleAPIKey string
		}{true, user, pc.Cfg.GoogleAPIKey})
		return
	}

	core.Error(w, http.StatusInternalServerError, "User not found")
	return
}

package controllers

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
	"se_ne/core"
	"se_ne/models"
)

type SessionController struct {
	Controller
}

// Login
func (sc *SessionController) Post(w http.ResponseWriter, r *http.Request) {
	var loginCredentials = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&loginCredentials)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(loginCredentials.Email) == 0 || len(loginCredentials.Password) == 0 {
		core.Error(w, http.StatusForbidden, "Access credentials have not correct values")
		return
	}

	user, err := models.GetUserByEmail(loginCredentials.Email)
	if err != nil {
		core.Error(w, http.StatusForbidden, "Invalid access credentials")
		return
	}
	ok := core.CheckPasswordHash(loginCredentials.Password, sc.Cfg.SecretKey, user.Password)
	if !ok {
		core.Error(w, http.StatusForbidden, "Invalid access credentials")
		return
	}
	token, err := Login(user, sc.Session, w, r)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.Data(w, struct {
		Token string `json:"token"`
	}{Token: token})
}

// Check that user is authenticated
func (sc *SessionController) Get(w http.ResponseWriter, r *http.Request) {
	session, err := sc.Session.Get(r, "user-session")
	if err != nil {
		core.Error(w, http.StatusInternalServerError, "Session not found")
		return
	}
	value, ok := session.Values["user-token"]
	valueStr, ok := value.(string)
	if ok && valueStr != "" {
		_, err := models.GetSessionByToken(valueStr)
		if err == nil {
			core.Data(w, struct {
				Status bool   `json:"status"`
				Token  string `json:"token"`
			}{true, valueStr})
			return
		}
	}

	core.Data(w, struct {
		Status bool   `json:"status"`
		Token  string `json:"token"`
	}{false, ""})
}

// Update session and generate new token
func (sc *SessionController) Put(w http.ResponseWriter, r *http.Request) {
	session, err := sc.Session.Get(r, "user-session")
	if err != nil {
		core.Error(w, http.StatusInternalServerError, "Session not found")
		return
	}
	value, ok := session.Values["user-token"]
	valueStr, ok := value.(string)
	if ok && valueStr != "" {
		user, err := models.GetUserByToken(valueStr)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "Token not found")
			return
		}
		models.DeleteSessionByToken(valueStr)
		token, err := Login(user, sc.Session, w, r)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "Token not found")
			return
		}

		core.Data(w, struct {
			Status   bool   `json:"status"`
			Redirect string `json:"redirect"`
			Token    string `json:"token"`
		}{true, "/", token})
		return
	} else {
		core.Error(w, http.StatusInternalServerError, "Token not found")
		return
	}
}

// Delete session and access token
func (sc *SessionController) Delete(w http.ResponseWriter, r *http.Request) {
	session, err := sc.Session.Get(r, "user-session")
	if err != nil {
		core.Error(w, http.StatusInternalServerError, "Session not found")
		return
	}
	value, ok := session.Values["user-token"]
	if ok && value != "" {
		models.DeleteSessionByToken(value.(string))
		delete(session.Values, "user-token")
		session.Save(r, w)
		core.Data(w, struct {
			Status   bool   `json:"status"`
			Redirect string `json:"redirect"`
		}{true, "/"})
		return
	} else {
		core.Error(w, http.StatusInternalServerError, "Token not found")
		return
	}
}

func (sc *SessionController) GoogleStartProcess(w http.ResponseWriter, r *http.Request) {
	oauthConfiguration := core.GenerateGoogleConfig(sc.Cfg)
	oauthStateString := "code" // TODO: generate token and write to session
	url := oauthConfiguration.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (sc *SessionController) GoogleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	oauthConfiguration := core.GenerateGoogleConfig(sc.Cfg)
	oauthStateString := "code"
	if state != oauthStateString {
		core.Error(w, http.StatusInternalServerError, "Token not found")
		return
	}
	code := r.FormValue("code") // TODO: read token from session
	token, err := oauthConfiguration.Exchange(oauth2.NoContext, code)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer response.Body.Close()
	cr := models.Credentials{}
	err = json.NewDecoder(response.Body).Decode(&cr)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, "could not handle user info")
		return
	}
	user, err := models.GetUserByEmail(cr.Email)
	if err != nil {
		//if  user don`t exist, try to do registration
		err = cr.Insert()
		if err != nil {
			core.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		user, err = models.GetUserByEmail(cr.Email)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Login user
	_, err = Login(user, sc.Session, w, r)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/", 200)
}

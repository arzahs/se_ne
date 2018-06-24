package controllers

import (
	"encoding/json"
	"net/http"
	"se_ne/core"
	"se_ne/models"
)

type UserController struct {
	Controller
}

// Create User
func (uc *UserController) Post(w http.ResponseWriter, r *http.Request) {
	var registrationForm = core.EmailRegistrationForm{}
	err := json.NewDecoder(r.Body).Decode(&registrationForm)
	if err != nil {
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validation block
	errorResponse := registrationForm.Validation()
	_, err = models.GetUserByEmail(registrationForm.Email)
	if err == nil {
		errorResponse = append(errorResponse, core.FormError{Message: "User is exist", Name: "email"})
	}
	if len(errorResponse) > 0 {
		core.Data(w, struct {
			Status    bool             `json:"status"`
			ErrorList []core.FormError `json:"errors"`
		}{false, errorResponse})
		return
	}

	// Create new user
	hashedPassword, _ := core.HashPassword(registrationForm.Password, uc.Cfg.SecretKey)
	userCredentials := models.Credentials{
		Email:    registrationForm.Email,
		Password: hashedPassword,
	}
	err = userCredentials.Insert()
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Login and redirect to profile
	user, err := models.GetUserByEmail(userCredentials.Email)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := Login(user, uc.Session, w, r)
	if err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.Data(w, struct {
		Status   bool   `json:"status"`
		Redirect string `json:"redirect"`
		Message  string `json:"message"`
		Token    string `json:"token"`
	}{true, "/", "Registration complete", token})
}

// Get info about user for profile page
func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if ok && userId != 0 {
		user, err := models.GetUserById(userId)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "User not found")
			return
		}
		core.Data(w, struct {
			Id        int    `json:"id"`
			Email     string `json:"email"`
			Address   string `json:"address"`
			Telephone string `json:"telephone"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}{user.Id, user.Email, user.Address, user.Telephone, user.LastName, user.FirstName})
		return
	}
	core.Error(w, http.StatusInternalServerError, "User not found")
	return
}

// Update user information
func (uc *UserController) Put(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if ok && userId != 0 {
		// check user
		user, err := models.GetUserById(userId)
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "User not found")
			return
		}

		// validation form
		profileForm := core.ProfileForm{}
		err = json.NewDecoder(r.Body).Decode(&profileForm)
		if err != nil {
			core.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		errorResponse := profileForm.Validation()
		if len(errorResponse) > 0 {
			core.Data(w, struct {
				Status    bool             `json:"status"`
				ErrorList []core.FormError `json:"errors"`
			}{false, errorResponse})
			return
		}

		// update structure
		user.Telephone = profileForm.Telephone
		user.Address = profileForm.Address
		user.FirstName = profileForm.FirstName
		user.LastName = profileForm.LastName
		user.Email = profileForm.Email
		user.IsActive = true // activate user after filled profile
		err = user.Update()
		if err != nil {
			core.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		core.Data(w, struct {
			Status   bool   `json:"status"`
			Redirect string `json:"redirect"`
			Message  string `json:"message"`
		}{true, "/", "Registration complete"})
		return
	}

	core.Error(w, http.StatusInternalServerError, "User not found")
	return
}

// Delete user information
func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if ok && userId != 0 {
		user, err := models.GetUserById(userId)
		if err != nil {
			core.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		session, err := uc.Session.Get(r, "user-session")
		if err != nil {
			core.Error(w, http.StatusInternalServerError, "Session not found")
			return
		}
		value, ok := session.Values["user-token"]
		if ok && value != "" {
			models.DeleteSessionByToken(value.(string))
			delete(session.Values, "user-token")
			session.Save(r, w)
			err = user.Remove()
			if err != nil {
				core.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
			core.Data(w, struct {
				Status bool `json:"status"`
			}{true})
			return
		}
	}
	core.Error(w, http.StatusInternalServerError, "User not found")
	return
}

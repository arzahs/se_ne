package controllers

import (
	"encoding/json"
	"se_ne/core"
	"net/http"
	"se_ne/models"
)

type ResetRequestController struct {
	Controller
}

// Generate request for password access
func (rrc *ResetRequestController) Post(w http.ResponseWriter, r *http.Request){
	var resetForm = struct {
		Email string   `json:"email"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&resetForm)
	if err != nil{
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	var errorResponse = make([]core.FormError, 0)
	if len(resetForm.Email) == 0{
		errorResponse = append(errorResponse, core.FormError{Message:"Email is not correct", Name: "email"})
	}

	user, err := models.GetUserByEmail(resetForm.Email)
	if err != nil{
		errorResponse = append(errorResponse, core.FormError{Message:"Email is not correct", Name: "email"})
	}
	if len(errorResponse) > 0{
		core.Data(w, struct {
			Status    bool             `json:"status"`
			ErrorList []core.FormError `json:"errors"`
		}{false, errorResponse})
		return
	}
	token, err := core.GenerateToken()
	if err != nil{
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	rs := models.ResetRequest{UserId:user.Id, Token: *token}
	err = rs.Insert()
	if err != nil{
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = core.SendResetPasswordToken(user.Email, *token, rrc.Cfg.Email)
	if err != nil{
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.Data(w, struct {
		Status bool `json:"status"`
	}{true})
	return
}

// Set New Password
func (rrc *ResetRequestController) Put(w http.ResponseWriter, r *http.Request){
	newPasswordForm := core.NewPasswordForm{}
	err := json.NewDecoder(r.Body).Decode(&newPasswordForm)
	if err != nil{
		core.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	var user = &models.User{}
	errorResponse := newPasswordForm.Validation()
	rps, err := models.GetResetRequestByToken(newPasswordForm.Token)
	if err != nil{
		errorResponse = append(errorResponse, core.FormError{Message:"Token is invalid", Name: "token"})
	}else {
		user, err = models.GetUserById(rps.UserId)
		if err != nil {
			errorResponse = append(errorResponse, core.FormError{Message: "Token is invalid", Name: "token"})
		}
	}
	if len(errorResponse) > 0 {
		core.Data(w, struct {
			Status    bool             `json:"status"`
			ErrorList []core.FormError `json:"errors"`
		}{false, errorResponse})
		return
	}
	hashedPassword, _ := core.HashPassword(newPasswordForm.Password, rrc.Cfg.SecretKey)
	userCredentials := models.Credentials{
		Email:    user.Email,
		Password: hashedPassword,
	}
	err = userCredentials.Update()
	if err != nil{
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	core.Data(w, struct {
		Status bool `json:"status"`
	}{true})
	return
}

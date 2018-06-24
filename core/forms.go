package core

import "regexp"

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type EmailRegistrationForm struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Password1 string `json:"password1"`
}

func (erf *EmailRegistrationForm) Validation() []FormError {
	var errorResponse = make([]FormError, 0)
	if len(erf.Email) == 0 || !rxEmail.MatchString(erf.Email){
		errorResponse = append(errorResponse, FormError{Message: "Email is not correct", Name: "email"})
	}

	if len(erf.Password) < 8 || len(erf.Password) > 38 {
		errorResponse = append(errorResponse, FormError{Message: "Password must have 8-38 characters", Name: "password"})
	}

	if erf.Password != erf.Password1 {
		errorResponse = append(errorResponse, FormError{Message: "Password and password confirmation must be equal", Name: "password1"})
	}
	return errorResponse
}

type ProfileForm struct {
	Email     string `json:"email"`
	Address   string `json:"address"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Telephone string `json:"telephone"`
}

func (pf *ProfileForm) Validation() []FormError {
	var errorResponse = make([]FormError, 0)
	if len(pf.Email) == 0 || !rxEmail.MatchString(pf.Email){
		errorResponse = append(errorResponse, FormError{Message: "Email is not correct", Name: "email"})
	}

	if len(pf.Address) == 0 {
		errorResponse = append(errorResponse, FormError{Message: "Address is not correct", Name: "address"})
	}

	if len(pf.FirstName) == 0 {
		errorResponse = append(errorResponse, FormError{Message: "First name is not correct", Name: "first_name"})
	}

	if len(pf.LastName) == 0 {
		errorResponse = append(errorResponse, FormError{Message: "Last name is not correct", Name: "last_name"})
	}

	if len(pf.Telephone) < 6 || len(pf.Telephone) > 12 {
		errorResponse = append(errorResponse, FormError{Message: "Telephone must have 6-12 characters", Name: "telephone"})
	}

	return errorResponse
}

type NewPasswordForm struct {
	Token     string `json:"token"`
	Password  string `json:"password"`
	Password1 string `json:"password1"`
}

func (npf *NewPasswordForm) Validation() []FormError {
	var errorResponse = make([]FormError, 0)

	if len(npf.Token) == 0 {
		errorResponse = append(errorResponse, FormError{Message: "Token is invalid", Name: "token"})
	}

	if len(npf.Password) < 8 || len(npf.Password) > 38 {
		errorResponse = append(errorResponse, FormError{Message: "Password must have 8-38 characters", Name: "password"})
	}

	if npf.Password != npf.Password1 {
		errorResponse = append(errorResponse, FormError{Message: "Password and password confirmation must be equal", Name: "password1"})
	}
	return errorResponse
}

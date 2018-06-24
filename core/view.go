package core

import (
	"html/template"
	"net/http"
)

//
//const(
//	HOME string = "profile.html"
//	LOGIN string = "login.html"
//	REGISTRATION string = "registration.html"
//	RESET_PASSWORD string = "forgot_password.html"
//	NEW_PASSWORD string = "new_password.html"
//)

type View struct{
	templates map[string]*template.Template
	pathToTemplate string
}

func NewView(pathToTemplate string) *View{
	templates := make(map[string]*template.Template)
	templates["profile"] =  template.Must(template.ParseFiles("./templates/profile.html", "./templates/base.html"))
	templates["login"] =  template.Must(template.ParseFiles("./templates/auth/login.html", "./templates/base.html"))
	templates["registration"] =  template.Must(template.ParseFiles("./templates/auth/registration.html", "./templates/base.html"))
	templates["reset"] =  template.Must(template.ParseFiles("./templates/auth/reset_password.html", "./templates/base.html"))
	templates["new_password"] =  template.Must(template.ParseFiles("./templates/auth/new_password.html", "./templates/base.html"))
	return &View{
		templates: templates,
		pathToTemplate: pathToTemplate,
	}
}

func (view *View) Render(w http.ResponseWriter, templateName string, context interface{}){
	tmpl, ok := view.templates[templateName]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, "base", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


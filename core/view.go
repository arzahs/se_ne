package core

import (
	"html/template"
	"net/http"
)

type View struct {
	templates       map[string]*template.Template
	pathToTemplates string
}

func NewView(pathToTemplates string) *View {
	templates := make(map[string]*template.Template)
	templates["profile"] = template.Must(template.ParseFiles(pathToTemplates+"profile.html", pathToTemplates+"base.html"))
	templates["login"] = template.Must(template.ParseFiles(pathToTemplates+"auth/login.html", pathToTemplates+"base.html"))
	templates["registration"] = template.Must(template.ParseFiles(pathToTemplates+"auth/registration.html", pathToTemplates+"base.html"))
	templates["reset"] = template.Must(template.ParseFiles(pathToTemplates+"auth/reset_password.html", pathToTemplates+"base.html"))
	templates["new_password"] = template.Must(template.ParseFiles(pathToTemplates+"auth/new_password.html", pathToTemplates+"base.html"))
	return &View{
		templates:       templates,
		pathToTemplates: pathToTemplates,
	}
}

func (view *View) Render(w http.ResponseWriter, templateName string, context interface{}) {
	tmpl, ok := view.templates[templateName]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, "base", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

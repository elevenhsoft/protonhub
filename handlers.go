package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = make(map[string]*template.Template)

func tmplCache() map[string]*template.Template {
	tmpl["index"] = template.Must(template.ParseFiles("./static/index.html", "./static/base.html"))

	tmpl["create"] = template.Must(template.ParseFiles("./static/create.html", "./static/base.html"))

	return tmpl
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmplCache()["index"].ExecuteTemplate(w, "base", nil)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	tmplCache()["create"].ExecuteTemplate(w, "base", nil)
}

func CreateDoneHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("launcherName")
	icon, handler, err := r.FormFile("launcherIcon")

	if err != nil {
		http.Error(w, "Error parsing file", http.StatusBadRequest)
		return
	}

	fmt.Println(name, icon, handler)
}

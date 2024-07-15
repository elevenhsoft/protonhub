package main

import (
	"html/template"
	"net/http"
	"strings"
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
	var launcherArgs []string

	name := r.FormValue("launcherName")
	proton := r.FormValue("protonPath")
	prefix := r.FormValue("prefixPath")
	exe := r.FormValue("launcherPathExe")
	gameId := r.FormValue("launcherGameId")
	store := r.FormValue("launcherGameStore")
	args := r.FormValue("launcherExeArgs")

	for _, split := range strings.Split(args, ",") {
		arg := strings.TrimSpace(split)
		launcherArgs = append(launcherArgs, arg)
	}

	obj := umu{
		Prefix:     prefix,
		Proton:     proton,
		GameID:     gameId,
		Exe:        exe,
		LaunchArgs: launcherArgs,
		Store:      store,
	}

	config_file := toTomlFileName(name)

	createTomlConfig(config_file, obj)
	conn := DbConnection()
	AddLauncherToDb(conn, name, args, obj)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var tmpl = make(map[string]*template.Template)

func tmplCache() map[string]*template.Template {
	editFunc := template.FuncMap{
		"unparseArgs": UnParseLauncherArgs,
	}

	editTpl, editErr := template.New("./static/edit.html").Funcs(editFunc).ParseFiles("./static/edit.html", "./static/base.html")

	tmpl["index"] = template.Must(template.ParseFiles("./static/index.html", "./static/base.html"))
	tmpl["edit"] = template.Must(editTpl, editErr)
	tmpl["create"] = template.Must(template.ParseFiles("./static/create.html", "./static/base.html"))

	return tmpl
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	conn := DbConnection()
	launchers := GetLaunchersFromDb(conn)

	tmplCache()["index"].ExecuteTemplate(w, "base", launchers)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	tmplCache()["create"].ExecuteTemplate(w, "base", nil)
}

func CreateDoneHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("launcherName")
	proton := r.FormValue("protonPath")
	prefix := r.FormValue("prefixPath")
	exe := r.FormValue("launcherPathExe")
	gameId := r.FormValue("launcherGameId")
	store := r.FormValue("launcherGameStore")
	args := r.FormValue("launcherExeArgs")

	obj := umu{
		Prefix:     prefix,
		Proton:     proton,
		GameID:     gameId,
		Exe:        exe,
		LaunchArgs: ParseLauncherArgs(args),
		Store:      store,
	}

	config_file := toTomlFileName(name)

	createTomlConfig(config_file, obj)
	conn := DbConnection()
	AddLauncherToDb(conn, config_file, name, args, obj)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	gameId := r.PathValue("gameId")

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, gameId).Config
	config := GetConfigPath(launcher)

	cmd := exec.Command("umu-run", "--config", config)

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	CmdToResponse(cmd, w)
}

type ProcessLockObject struct {
	Pid    string `json:"pid"`
	GameID string `json:"gameid"`
}

func CreateProcessLockHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var obj ProcessLockObject
	err = json.Unmarshal(body, &obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	CreateLockfileForProcess(obj.GameID, obj.Pid)
}

func RemoveProcessLockHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var obj ProcessLockObject
	err = json.Unmarshal(body, &obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RemoveLockfileForProcess(obj.GameID)
}

func StopProcessHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("gameId")

	if KillProcessForGameId(gameId) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RunWinetricksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	gameId := r.PathValue("gameId")
	verbs_get := r.PathValue("verbs")

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, gameId)

	var verbs []string
	verbs = append(verbs, "winetricks")
	for _, arg := range strings.Split(verbs_get, " ") {
		verbs = append(verbs, arg)
	}

	gameidEnv := fmt.Sprintf("GAMEID=%s", launcher.GameID)
	protonpathEnv := fmt.Sprintf("PROTONPATH=%s", launcher.Proton)

	cmd := exec.Command("umu-run", verbs...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, gameidEnv, protonpathEnv)

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	CmdToResponse(cmd, w)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("gameId")

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, gameId)

	tmplCache()["edit"].ExecuteTemplate(w, "base", launcher)
}

func EditDoneHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("launcherName")
	proton := r.FormValue("protonPath")
	prefix := r.FormValue("prefixPath")
	exe := r.FormValue("launcherPathExe")
	gameId := r.FormValue("launcherGameId")
	store := r.FormValue("launcherGameStore")
	args := r.FormValue("launcherExeArgs")

	obj := umu{
		Prefix:     prefix,
		Proton:     proton,
		GameID:     gameId,
		Exe:        exe,
		LaunchArgs: ParseLauncherArgs(args),
		Store:      store,
	}

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, gameId)

	conn = DbConnection()
	UpdateLauncherInDb(conn, launcher.Config, name, args, obj)
	updateTomlFile(launcher.Config, obj)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

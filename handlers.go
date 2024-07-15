package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
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

type RunHandlerObject struct {
	GameID string `json:"gameId"`
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var obj RunHandlerObject
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Fatal(err)
	}

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, obj.GameID).Config
	config := GetConfigPath(launcher)

	cmd := exec.Command("umu-run", "--config", config)
	_, err = cmd.Output()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

type RunWinetricksHandlerObject struct {
	GameID string `json:"gameId"`
	Verbs  string `json:"verbs"`
}

type WinetricksLogResponse struct {
	Log string `json:"log"`
}

func RunWinetricksHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var obj RunWinetricksHandlerObject
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Fatal(err)
	}

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, obj.GameID)

	var verbs []string
	verbs = append(verbs, "winetricks")
	for _, arg := range strings.Split(obj.Verbs, " ") {
		verbs = append(verbs, arg)
	}

	gameidEnv := fmt.Sprintf("GAMEID=%s", launcher.GameID)
	protonpathEnv := fmt.Sprintf("PROTONPATH=%s", launcher.Proton)

	cmd := exec.Command("umu-run", verbs...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, gameidEnv, protonpathEnv)

	out, _ := cmd.CombinedOutput()

	response := WinetricksLogResponse{Log: fmt.Sprintf("%s", string(out))}
	resp_log, err := json.Marshal(response)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp_log)
	w.(http.Flusher).Flush()
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

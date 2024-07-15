package main

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func DbConnection() *sql.DB {
	location := filepath.Join(phStorePath(), "database.db")
	db, err := sql.Open("sqlite3", location)

	if err != nil {
		log.Fatal("Can't open database")
	}

	return db
}

func LaunchersTableInit(conn *sql.DB) {
	sql := `
CREATE TABLE IF NOT EXISTS launchers (
	config text not null primary key,
	name text,
	prefix text,
	proton text,
	game_id text,
	exefile text,
	args text,
	store text
);`
	conn.Exec(sql)
	defer conn.Close()
}

func AddLauncherToDb(conn *sql.DB, config, name, args string, launcher umu) {
	tx, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql := `
INSERT INTO launchers(config, name, prefix, proton, game_id, exefile, args, store) values (?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	stmt.Exec(config, name, launcher.Prefix, launcher.Proton, launcher.GameID, launcher.Exe, args, launcher.Store)
	tx.Commit()
}

func GetLaunchersFromDb(conn *sql.DB) []Launcher {
	sql := `
SELECT config, name, prefix, proton, game_id, exefile, args, store FROM launchers 
	`
	rows, err := conn.Query(sql)

	if err != nil {
		log.Fatal(err)
	}

	var launchers []Launcher

	for rows.Next() {
		var config string
		var name string
		var prefix string
		var proton string
		var game_id string
		var exefile string
		var args string
		var store string

		err = rows.Scan(&config, &name, &prefix, &proton, &game_id, &exefile, &args, &store)
		if err != nil {
			log.Fatal(err)
		}

		launchers = append(launchers, Launcher{
			Config:     config,
			Name:       name,
			Prefix:     prefix,
			Proton:     proton,
			GameID:     game_id,
			Exe:        exefile,
			LaunchArgs: ParseLauncherArgs(args),
			Store:      store,
		})
	}

	return launchers
}

func GetLauncherByIdFromDb(conn *sql.DB, gameId string) Launcher {
	sql := `
SELECT config, name, prefix, proton, game_id, exefile, args, store FROM launchers WHERE game_id = ?
	`
	rows, err := conn.Query(sql, gameId)

	if err != nil {
		log.Fatal(err)
	}

	var launcher Launcher

	for rows.Next() {
		var config string
		var name string
		var prefix string
		var proton string
		var game_id string
		var exefile string
		var args string
		var store string

		err = rows.Scan(&config, &name, &prefix, &proton, &game_id, &exefile, &args, &store)
		if err != nil {
			log.Fatal(err)
		}

		launcher = Launcher{
			Config:     config,
			Name:       name,
			Prefix:     prefix,
			Proton:     proton,
			GameID:     game_id,
			Exe:        exefile,
			LaunchArgs: ParseLauncherArgs(args),
			Store:      store,
		}
	}

	return launcher
}

func UpdateLauncherInDb(conn *sql.DB, config, name, args string, launcher umu) {
	sql := `
UPDATE launchers SET name = ?, prefix = ?, proton = ?, game_id = ?, exefile = ?, args = ?, store = ?  WHERE config = ?
	`
	conn.Exec(sql, name, launcher.Prefix, launcher.Proton, launcher.GameID, launcher.Exe, args, launcher.Store, config)
	defer conn.Close()
}

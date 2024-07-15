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
	id integer not null primary key,
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

func AddLauncherToDb(conn *sql.DB, name, args string, launcher umu) {
	tx, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql := `
INSERT INTO launchers(name, prefix, proton, game_id, exefile, args, store) values (?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	stmt.Exec(name, launcher.Prefix, launcher.Proton, launcher.GameID, launcher.Exe, args, launcher.Store)
	tx.Commit()
}

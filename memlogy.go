package main

import (
	"memlogy/cmd"
	"memlogy/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.OpenDatabase()
	cmd.Execute()
}

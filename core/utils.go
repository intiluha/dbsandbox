package core

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	Driver       = "mysql"
	Username     = "root"
	Password     = "pass123"
	Hostname     = "127.0.0.1:3306"
	DatabaseName = "queue"
	Table        = "Queue"
	IdField      = "id"
	DataField    = "data"
	Order        = "ASC"  // Change to "DESC" for LIFO
)

func dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, db)
}

func databaseCloser(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("error [%s] when closing DB\n", err)
	}
}
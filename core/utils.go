package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	Driver       = "mysql"
	Hostname     = "127.0.0.1:3306"
	DatabaseName = "queue"
	Table        = "Queue"
	IdField      = "id"
	DataField    = "data"
	Order        = "ASC" // Change to "DESC" for LIFO
)

var (
	Username = os.Getenv("MYSQL_USER")
	Password = os.Getenv("MYSQL_PASS")
)

type Row struct {
	ID   uint
	Data string
}

func (Row) TableName() string {
	return Table
}

func dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, db)
}

func databaseCloser(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("%s when closing DB\n", err)
	}
}

func Assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

package core

import (
	"database/sql"
	"fmt"
	"log"
)

type orderType bool

const (
	AscendingOrder  = orderType(false)
	DescendingOrder = orderType(true)
)

type row struct {
	ID   uint
	Data string
}

func (row) TableName() string {
	return table
}

func dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, db)
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

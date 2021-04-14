package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "pass123"
	hostname = "127.0.0.1:3306"
	dbname   = "queue"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createDB(dbName string) error {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer db.Close()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelFunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("error [%s] when creating DB\n", err)
	}
	no, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error [%s] when fetching rows", err)
	}
	log.Printf("rows affected %d\n", no)
	return nil
}

func main() {
	err := createDB(dbname)
	if err != nil {
		log.Fatal(err)
	}
}
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
)

func dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, db)
}

func createDatabaseSQL() string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", DatabaseName)
}

func createTableSQL() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY AUTO_INCREMENT, %s VARCHAR(255) NOT NULL);", Table, IdField, DataField)
}

func insertSQL(data string) string {
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", Table, DataField, data)
}

func createTable() error {
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error [%s] when closing DB\n", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, createTableSQL())
	if err != nil {
		return fmt.Errorf("error [%s] when creating table\n", err)
	}
	no, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error [%s] when fetching rows", err)
	}
	log.Printf("rows affected %d\n", no)
	return nil
}

func createDatabase() error {
	db, err := sql.Open(Driver, dsn(""))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error [%s] when closing DB\n", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, createDatabaseSQL())
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

func writer(pattern string, n int) error {
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error [%s] when closing DB\n", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n) * time.Second)
	defer cancel()
	for i := 0; i < n; i++ {
		_, err := db.ExecContext(ctx, insertSQL(strconv.Itoa(i)))
		if err != nil {
			return fmt.Errorf("error [%s] when creating table\n", err)
		}
	}
	return nil
}

func main() {
	err := writer("", 2000)
	if err != nil {
		log.Fatal(err)
	}
}
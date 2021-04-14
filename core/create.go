package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func CreateDatabase() error {
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(""))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func CreateTable() error {
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		return fmt.Errorf("error [%s] when opening DB\n", err)
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

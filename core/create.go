package core

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func CreateDatabase() error {
	// Open database, defer closing
	db, err := sql.Open(driver, dsn(""))
	if err != nil {
		return fmt.Errorf("%s in CreateDatabase when opening DB", err)
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, createDatabaseSQL())
	if err != nil {
		return fmt.Errorf("%s in CreateDatabase when creating DB", err)
	}
	return nil
}

func CreateTable() error {
	// Open database, defer closing
	db, err := sql.Open(driver, dsn(databaseName))
	if err != nil {
		return fmt.Errorf("%s in CreateTable when opening DB", err)
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, createTableSQL())
	if err != nil {
		return fmt.Errorf("%s in CreateTable when creating table", err)
	}
	return nil
}

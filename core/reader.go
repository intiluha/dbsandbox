package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func process(data string) {
	fmt.Println(data)
}

func Reader(n int, errs chan error) {
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		errs <- fmt.Errorf("error [%s] when opening DB\n", err)
		return
	}
	defer databaseCloser(db)

	// Create context and establish connection, defer canceling and closing
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	conn, err := db.Conn(ctx)
	if err != nil {
		cancel()
		errs <- fmt.Errorf("error [%s] when establishing connection\n", err)
		return
	}
	defer func() {
		cancel()
		_ = conn.Close()
	}()

	var id int
	var data string
	for i := 0; i < n; i++ {
		// Begin transaction
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			errs <- fmt.Errorf("error [%s] when starting transaction\n", err)
			return
		}
		// Try to scan first row
		err = tx.QueryRow(queryOneSQL()).Scan(&id, &data)
		if err != nil {
			if err != sql.ErrNoRows {
				errs <- fmt.Errorf("error [%s] when creating table\n", err)
				return
			}
			// If no rows are present, wait
			_ = tx.Commit()
			time.Sleep(time.Second)
			i--
			log.Print("sleep")
		} else {
			_, err := tx.Exec(deleteSQL(id))
			if err != nil {
				// If row was already deleted, rollback
				err = tx.Rollback()
			} else {
				// If everything is Okay, commit and process
				err = tx.Commit()
				process(data)
			}
			if err != nil {
				errs <- fmt.Errorf("error [%s] when comitting/rolling back transaction\n", err)
				return
			}
		}
	}
	errs <- nil
	return
}

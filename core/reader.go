package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func process(data string, n int64) {
	fmt.Println(data, n)
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
				errs <- fmt.Errorf("error [%s] when querying row\n", err)
				return
			}
			// If no rows are present, wait
			_ = tx.Commit()
			time.Sleep(time.Second)
			i--
			log.Print("sleep")
			continue
		}
		// Delete row
		res, err := tx.Exec(deleteSQL(id))
		if err != nil {
			errs <- fmt.Errorf("error [%s] when deleting row\n", err)
			return
		}
		nRows, err := res.RowsAffected()
		if err != nil {
			errs <- fmt.Errorf("error [%s] when fetching rows\n", err)
			return
		}
		if nRows == 0 {
			// If row was already deleted by another reader, rollback
			log.Print("rollback")
			err = tx.Rollback()
			i--
		} else {
			// If everything is Okay, commit and process
			err = tx.Commit()
			process(data, nRows)
		}
		if err != nil {
			errs <- fmt.Errorf("error [%s] when comitting/rolling back transaction\n", err)
			return
		}
	}
	errs <- nil
	return
}

package core

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func generateSequence(prefix string, n int) []string {
	slice := make([]string, n)
	for i := range slice {
		slice[i] = prefix+strconv.Itoa(i)
	}
	return slice
}

func Writer(prefix string, n int, errs chan error) {
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		errs <- fmt.Errorf("error [%s] when opening DB\n", err)
		return
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	defer cancel()
	for _, x := range generateSequence(prefix, n){
		_, err := db.ExecContext(ctx, insertSQL(x))
		if err != nil {
			errs <- fmt.Errorf("error [%s] inserting row\n", err)
			return
		}
	}
	errs <- nil
	return
}

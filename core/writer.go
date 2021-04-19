package core

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func generateSequence(prefix string, n int) []string {
	slice := make([]string, n)
	for i := range slice {
		slice[i] = prefix + strconv.Itoa(i)
	}
	return slice
}

func WriterORM(prefix string, n int, errs chan error) {
	db, err := gorm.Open(mysql.Open(dsn(DatabaseName)), &gorm.Config{})
	if err != nil {
		errs <- err
		return
	}

	for _, x := range generateSequence(prefix, n) {
		err = db.Create(&Row{Data: x}).Error
		if err != nil {
			errs <- err
			return
		}
	}
	errs <- nil
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
	for _, x := range generateSequence(prefix, n) {
		_, err := db.ExecContext(ctx, insertSQL(x))
		if err != nil {
			errs <- fmt.Errorf("error [%s] inserting row\n", err)
			return
		}
	}
	errs <- nil
}

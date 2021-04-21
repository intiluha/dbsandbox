package core

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

func generateSequence(prefix string, n int) []string {
	slice := make([]string, n)
	for i := range slice {
		slice[i] = prefix + strconv.Itoa(i)
	}
	return slice
}

func WriterORM(prefix string, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := gorm.Open(mysql.Open(dsn(databaseName)), &gorm.Config{})
	if err != nil {
		log.Println(err, "in WriterORM when opening db")
		return
	}

	for _, x := range generateSequence(prefix, n) {
		err = db.Create(&row{Data: x}).Error
		if err != nil {
			log.Println(err, "in WriterORM when inserting element")
			return
		}
	}
}

func Writer(prefix string, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database, defer closing
	db, err := sql.Open(driver, dsn(databaseName))
	if err != nil {
		log.Println(err, "in Writer when opening DB")
		return
	}
	defer databaseCloser(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	defer cancel()
	for _, x := range generateSequence(prefix, n) {
		_, err := db.ExecContext(ctx, insertSQL(), x)
		if err != nil {
			log.Println(err, "in Writer when inserting row")
			return
		}
	}
}

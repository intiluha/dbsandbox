package core

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sync"
	"time"
)

func process(data string) {
	fmt.Println(data)
}

func interrupt() {
	log.Print("interrupted")
}

func sleep() {
	time.Sleep(10*time.Millisecond)
	log.Print("sleep")
}

func ReaderORM(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database
	db, err := gorm.Open(mysql.Open(dsn(DatabaseName)), &gorm.Config{})
	if err != nil {
		log.Println(err, "in ReaderORM when opening DB")
		return
	}

	var row Row
	for i := 0; i < n; {
		row = Row{}
		err = db.First(&row).Error
		if err == logger.ErrRecordNotFound {
			sleep()
			continue
		}
		if err != nil {
			log.Println(err, "in ReaderORM querying element")
			return
		}
		res := db.Delete(row)
		if err = res.Error; err != nil {
			log.Println(err, "in ReaderORM when opening DB")
			return
		}
		if res.RowsAffected == 0 {
			interrupt()
		} else {
			process(row.Data)
			i++
		}
	}
}

func Reader(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		log.Println(err, "in Reader when opening DB")
		return
	}
	defer databaseCloser(db)

	// Create context and establish connection, defer canceling and closing
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	conn, err := db.Conn(ctx)
	if err != nil {
		cancel()
		log.Println(err, "in Reader when establishing connection")
		return
	}
	defer func() {
		cancel()
		_ = conn.Close()
	}()

	var id int
	var data string
	for i := 0; i < n; {
		// Begin transaction TODO: tx not needed
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			log.Println(err, "in Reader when beginning transaction")
			return
		}
		// Try to scan first row
		err = tx.QueryRow(queryOneSQL()).Scan(&id, &data)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Println(err, "in Reader when querying row")
				return
			}
			// If no rows are present, wait
			_ = tx.Commit()
			sleep()
			continue
		}
		// Delete row
		res, err := tx.Exec(deleteSQL(id))
		if err != nil {
			log.Println(err, "in Reader when deleting row")
			return
		}
		nRows, err := res.RowsAffected()
		if err != nil {
			log.Println(err, "in Reader when fetching rows")
			return
		}
		if nRows == 0 {
			// If row was already deleted by another reader, rollback
			interrupt()
			err = tx.Rollback()
		} else {
			// If everything is Okay, commit and process
			err = tx.Commit()
			process(data)
			i++
		}
		if err != nil {
			log.Println(err, "in Reader when committing/rolling back transaction")
			return
		}
	}
}

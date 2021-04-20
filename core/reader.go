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

func ReaderORM(order orderType, n int, wg *sync.WaitGroup) {
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
		if order == AscendingOrder {
			err = db.First(&row).Error
		} else {
			err = db.Last(&row).Error
		}

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

func Reader(order orderType, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database, defer closing
	db, err := sql.Open(Driver, dsn(DatabaseName))
	if err != nil {
		log.Println(err, "in Reader when opening DB")
		return
	}
	defer databaseCloser(db)

	// Create context, defer canceling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	defer cancel()

	var id int
	var data string
	for i := 0; i < n; {
		// Try to scan first row
		err = db.QueryRowContext(ctx, queryOneSQL(order)).Scan(&id, &data)
		if err == sql.ErrNoRows {
			sleep()
			continue
		} else if err != nil {
			log.Println(err, "in Reader when querying row")
			return
		}
		// Delete row
		res, err := db.ExecContext(ctx, deleteSQL(id))
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
		} else {
			// If everything is Okay, commit and process
			process(data)
			i++
		}
	}
}

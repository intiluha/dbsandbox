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
	time.Sleep(10 * time.Millisecond)
	log.Print("sleep")
}

func ReaderORM(order orderType, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database
	db, err := gorm.Open(mysql.Open(dsn(databaseName)), &gorm.Config{})
	if err != nil {
		log.Println(err, "in ReaderORM when opening DB")
		return
	}

	var r row
	for i := 0; i < n; {
		r = row{}
		if order == AscendingOrder {
			err = db.First(&r).Error
		} else {
			err = db.Last(&r).Error
		}

		if err == logger.ErrRecordNotFound {
			sleep()
			continue
		}
		if err != nil {
			log.Println(err, "in ReaderORM querying element")
			return
		}
		res := db.Delete(r)
		if err = res.Error; err != nil {
			log.Println(err, "in ReaderORM when opening DB")
			return
		}
		if res.RowsAffected == 0 {
			interrupt()
		} else {
			process(r.Data)
			i++
		}
	}
}

func Reader(order orderType, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Open database, defer closing
	db, err := sql.Open(driver, dsn(databaseName))
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
		res, err := db.ExecContext(ctx, deleteSQL(), id)
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

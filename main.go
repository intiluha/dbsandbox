package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/intiluha/dbsandbox/core"
	"log"
)

func main() {
	err := core.CreateDatabase()
	core.Assert(err)
	err = core.CreateTable()
	core.Assert(err)

	nWriters, nReaders, nOperations := 2, 2, 200
	errs := make(chan error, nWriters+nReaders)
	for i := 0; i < nWriters; i++ {
		go core.WriterORM(string(rune('a'+i)), nOperations, errs)
	}
	for i := 0; i < nReaders; i++ {
		go core.ReaderORM(nOperations, errs)
	}
	for i := 0; i < nWriters+nReaders; i++ {
		log.Print(<-errs)
	}
}
